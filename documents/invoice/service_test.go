// +build unit

package invoice

import (
	"context"
	"fmt"
	"math/big"
	"strconv"
	"testing"

	"github.com/centrifuge/go-centrifuge/testingutils"
	"github.com/centrifuge/go-centrifuge/testingutils/commons"
	"github.com/centrifuge/go-centrifuge/testingutils/documents"

	"github.com/centrifuge/centrifuge-protobufs/gen/go/coredocument"
	"github.com/centrifuge/go-centrifuge/anchors"
	"github.com/centrifuge/go-centrifuge/code"
	"github.com/centrifuge/go-centrifuge/config"
	"github.com/centrifuge/go-centrifuge/coredocument"
	"github.com/centrifuge/go-centrifuge/documents"
	"github.com/centrifuge/go-centrifuge/identity"
	clientinvoicepb "github.com/centrifuge/go-centrifuge/protobufs/gen/go/invoice"
	"github.com/centrifuge/go-centrifuge/signatures"
	"github.com/centrifuge/go-centrifuge/utils"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	centID  = utils.RandomSlice(identity.CentIDLength)
	key1Pub = [...]byte{230, 49, 10, 12, 200, 149, 43, 184, 145, 87, 163, 252, 114, 31, 91, 163, 24, 237, 36, 51, 165, 8, 34, 104, 97, 49, 114, 85, 255, 15, 195, 199}
	key1    = []byte{102, 109, 71, 239, 130, 229, 128, 189, 37, 96, 223, 5, 189, 91, 210, 47, 89, 4, 165, 6, 188, 53, 49, 250, 109, 151, 234, 139, 57, 205, 231, 253, 230, 49, 10, 12, 200, 149, 43, 184, 145, 87, 163, 252, 114, 31, 91, 163, 24, 237, 36, 51, 165, 8, 34, 104, 97, 49, 114, 85, 255, 15, 195, 199}
)

type mockAnchorRepo struct {
	mock.Mock
	anchors.AnchorRepository
}

func (r *mockAnchorRepo) GetDocumentRootOf(anchorID anchors.AnchorID) (anchors.DocRoot, error) {
	args := r.Called(anchorID)
	docRoot, _ := args.Get(0).(anchors.DocRoot)
	return docRoot, args.Error(1)
}

func TestDefaultService(t *testing.T) {
	srv := DefaultService(getRepository(), &testingutils.MockCoreDocumentProcessor{}, nil)
	assert.NotNil(t, srv, "must be non-nil")
}

func getServiceWithMockedLayers() Service {
	return DefaultService(getRepository(), &testingutils.MockCoreDocumentProcessor{}, &mockAnchorRepo{})
}

func createMockDocument() (*Invoice, error) {
	documentIdentifier := utils.RandomSlice(32)
	nextIdentifier := utils.RandomSlice(32)
	inv1 := &Invoice{
		InvoiceNumber: "test_invoice",
		GrossAmount:   60,
		CoreDocument: &coredocumentpb.CoreDocument{
			DocumentIdentifier: documentIdentifier,
			CurrentVersion:     documentIdentifier,
			NextVersion:        nextIdentifier,
		},
	}
	err := getRepository().Create(documentIdentifier, inv1)
	return inv1, err
}

func TestService_DeriveFromCoreDocument(t *testing.T) {
	// nil doc
	invSrv := service{repo: getRepository()}
	_, err := invSrv.DeriveFromCoreDocument(nil)
	assert.Error(t, err, "must fail to derive")

	// successful
	data := testingdocuments.CreateInvoiceData()
	coreDoc := testingdocuments.CreateCDWithEmbeddedInvoice(t, data)
	model, err := invSrv.DeriveFromCoreDocument(coreDoc)
	assert.Nil(t, err, "must return model")
	assert.NotNil(t, model, "model must be non-nil")
	inv, ok := model.(*Invoice)
	assert.True(t, ok, "must be true")
	assert.Equal(t, inv.Payee[:], data.Payee)
	assert.Equal(t, inv.Sender[:], data.Sender)
	assert.Equal(t, inv.Recipient[:], data.Recipient)
	assert.Equal(t, inv.GrossAmount, data.GrossAmount)
}

func TestService_DeriveFromPayload(t *testing.T) {
	invSrv := getServiceWithMockedLayers()
	payload := testingdocuments.CreateInvoicePayload()
	var model documents.Model
	var err error

	// fail due to nil payload
	_, err = invSrv.DeriveFromCreatePayload(nil, nil)
	assert.Error(t, err, "DeriveWithInvoiceInput should produce an error if invoiceInput equals nil")

	contextHeader, err := documents.NewContextHeader()
	assert.Nil(t, err)
	model, err = invSrv.DeriveFromCreatePayload(payload, contextHeader)
	assert.Nil(t, err, "valid invoiceData shouldn't produce an error")

	receivedCoreDocument, err := model.PackCoreDocument()
	assert.Nil(t, err, "model should be able to return the core document with embedded invoice")
	assert.NotNil(t, receivedCoreDocument.EmbeddedData, "embeddedData should be field")
}

func TestService_GetLastVersion(t *testing.T) {
	invSrv := getServiceWithMockedLayers()
	thirdIdentifier := utils.RandomSlice(32)
	doc, err := createMockDocument()
	assert.Nil(t, err)

	mod1, err := invSrv.GetCurrentVersion(doc.CoreDocument.DocumentIdentifier)
	assert.Nil(t, err)

	invLoad1, _ := mod1.(*Invoice)
	assert.Equal(t, invLoad1.CoreDocument.CurrentVersion, doc.CoreDocument.DocumentIdentifier)

	inv2 := &Invoice{
		GrossAmount: 60,
		CoreDocument: &coredocumentpb.CoreDocument{
			DocumentIdentifier: doc.CoreDocument.DocumentIdentifier,
			CurrentVersion:     doc.CoreDocument.NextVersion,
			NextVersion:        thirdIdentifier,
		},
	}

	err = getRepository().Create(doc.CoreDocument.NextVersion, inv2)
	assert.Nil(t, err)

	mod2, err := invSrv.GetCurrentVersion(doc.CoreDocument.DocumentIdentifier)
	assert.Nil(t, err)

	invLoad2, _ := mod2.(*Invoice)
	assert.Equal(t, invLoad2.CoreDocument.CurrentVersion, doc.CoreDocument.NextVersion)
	assert.Equal(t, invLoad2.CoreDocument.NextVersion, thirdIdentifier)
}

func TestService_GetVersion_invalid_version(t *testing.T) {
	invSrv := getServiceWithMockedLayers()
	currentVersion := utils.RandomSlice(32)
	inv := &Invoice{
		GrossAmount: 60,
		CoreDocument: &coredocumentpb.CoreDocument{
			DocumentIdentifier: utils.RandomSlice(32),
			CurrentVersion:     currentVersion,
		},
	}
	err := getRepository().Create(currentVersion, inv)
	assert.Nil(t, err)

	mod, err := invSrv.GetVersion(utils.RandomSlice(32), currentVersion)
	assert.EqualError(t, err, "[4]document not found for the given version: version is not valid for this identifier")
	assert.Nil(t, mod)
}

func TestService_GetVersion(t *testing.T) {
	invSrv := getServiceWithMockedLayers()
	documentIdentifier := utils.RandomSlice(32)
	currentVersion := utils.RandomSlice(32)
	inv := &Invoice{
		GrossAmount: 60,
		CoreDocument: &coredocumentpb.CoreDocument{
			DocumentIdentifier: documentIdentifier,
			CurrentVersion:     currentVersion,
		},
	}
	err := getRepository().Create(currentVersion, inv)
	assert.Nil(t, err)

	mod, err := invSrv.GetVersion(documentIdentifier, currentVersion)
	assert.Nil(t, err)
	loadInv, _ := mod.(*Invoice)
	assert.Equal(t, loadInv.CoreDocument.CurrentVersion, currentVersion)
	assert.Equal(t, loadInv.CoreDocument.DocumentIdentifier, documentIdentifier)

	mod, err = invSrv.GetVersion(documentIdentifier, []byte{})
	assert.Error(t, err)
}

func TestService_Create(t *testing.T) {
	ctxh, err := documents.NewContextHeader()
	assert.Nil(t, err)
	invSrv := service{repo: getRepository()}
	ctx := context.Background()

	// calculate data root fails
	m, err := invSrv.Create(context.Background(), &testingdocuments.MockModel{})
	assert.Nil(t, m)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown document type")

	// anchor fails
	po, err := invSrv.DeriveFromCreatePayload(testingdocuments.CreateInvoicePayload(), ctxh)
	assert.Nil(t, err)
	proc := &testingutils.MockCoreDocumentProcessor{}
	proc.On("PrepareForSignatureRequests", po).Return(fmt.Errorf("anchoring failed")).Once()
	invSrv.coreDocProcessor = proc
	m, err = invSrv.Create(ctx, po)
	proc.AssertExpectations(t)
	assert.Nil(t, m)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "anchoring failed")

	// success
	po, err = invSrv.DeriveFromCreatePayload(testingdocuments.CreateInvoicePayload(), ctxh)
	assert.Nil(t, err)
	proc = &testingutils.MockCoreDocumentProcessor{}
	proc.On("PrepareForSignatureRequests", po).Return(nil).Once()
	proc.On("RequestSignatures", ctx, po).Return(nil).Once()
	proc.On("PrepareForAnchoring", po).Return(nil).Once()
	proc.On("AnchorDocument", po).Return(nil).Once()
	proc.On("SendDocument", ctx, po).Return(nil).Once()
	invSrv.coreDocProcessor = proc
	m, err = invSrv.Create(ctx, po)
	proc.AssertExpectations(t)
	assert.Nil(t, err)

	newCD, err := m.PackCoreDocument()
	assert.Nil(t, err)
	assert.True(t, getRepository().Exists(newCD.DocumentIdentifier))
	assert.True(t, getRepository().Exists(newCD.CurrentVersion))
}

func TestService_DeriveInvoiceData(t *testing.T) {
	invSrv := getServiceWithMockedLayers()

	// some random model
	_, err := invSrv.DeriveInvoiceData(&mockModel{})
	assert.Error(t, err, "Derive must fail")

	// success
	payload := testingdocuments.CreateInvoicePayload()
	contextHeader, err := documents.NewContextHeader()
	assert.Nil(t, err)
	inv, err := invSrv.DeriveFromCreatePayload(payload, contextHeader)
	assert.Nil(t, err, "must be non nil")
	data, err := invSrv.DeriveInvoiceData(inv)
	assert.Nil(t, err, "Derive must succeed")
	assert.NotNil(t, data, "data must be non nil")
}

func TestService_DeriveInvoiceResponse(t *testing.T) {
	// success
	invSrv := service{repo: getRepository()}
	payload := testingdocuments.CreateInvoicePayload()
	contextHeader, err := documents.NewContextHeader()
	assert.Nil(t, err)
	inv1, err := invSrv.DeriveFromCreatePayload(payload, contextHeader)
	assert.Nil(t, err, "must be non nil")
	inv, ok := inv1.(*Invoice)
	assert.True(t, ok)
	inv.CoreDocument = &coredocumentpb.CoreDocument{
		DocumentIdentifier: []byte{},
	}
	resp, err := invSrv.DeriveInvoiceResponse(inv)
	assert.Nil(t, err, "Derive must succeed")
	assert.NotNil(t, resp, "data must be non nil")
	assert.Equal(t, resp.Data, payload.Data, "data mismatch")
}

// Functions returns service mocks
func mockSignatureCheck(i *Invoice, invSrv Service) identity.Service {
	idkey := &identity.EthereumIdentityKey{
		Key:       key1Pub,
		Purposes:  []*big.Int{big.NewInt(identity.KeyPurposeSigning)},
		RevokedAt: big.NewInt(0),
	}
	anchorID, _ := anchors.NewAnchorID(i.CoreDocument.DocumentIdentifier)
	docRoot, _ := anchors.NewDocRoot(i.CoreDocument.DocumentRoot)
	mockRepo := invSrv.(service).anchorRepository.(*mockAnchorRepo)
	mockRepo.On("GetDocumentRootOf", anchorID).Return(docRoot, nil).Once()
	srv := &testingcommons.MockIDService{}
	id := &testingcommons.MockID{}
	centID, _ := identity.ToCentID(centID)
	srv.On("LookupIdentityForID", centID).Return(id, nil).Once()
	id.On("FetchKey", key1Pub[:]).Return(idkey, nil).Once()
	return srv
}

func setIdentityService(idService identity.Service) {
	identity.IDService = idService
}

func TestService_CreateProofs(t *testing.T) {
	defer setIdentityService(identity.IDService)
	invSrv := getServiceWithMockedLayers()
	i, err := createAnchoredMockDocument(t, false)
	assert.Nil(t, err)
	idService := mockSignatureCheck(i, invSrv)
	setIdentityService(idService)
	proof, err := invSrv.CreateProofs(i.CoreDocument.DocumentIdentifier, []string{"invoice_number"})
	assert.Nil(t, err)
	assert.Equal(t, i.CoreDocument.DocumentIdentifier, proof.DocumentId)
	assert.Equal(t, i.CoreDocument.DocumentIdentifier, proof.VersionId)
	assert.Equal(t, len(proof.FieldProofs), 1)
	assert.Equal(t, proof.FieldProofs[0].GetProperty(), "invoice_number")
}

func TestService_CreateProofsValidationFails(t *testing.T) {
	defer setIdentityService(identity.IDService)
	invSrv := getServiceWithMockedLayers()
	i, err := createAnchoredMockDocument(t, false)
	assert.Nil(t, err)
	i.CoreDocument.SigningRoot = nil
	err = getRepository().Update(i.CoreDocument.CurrentVersion, i)
	assert.Nil(t, err)
	idService := mockSignatureCheck(i, invSrv)
	setIdentityService(idService)
	_, err = invSrv.CreateProofs(i.CoreDocument.DocumentIdentifier, []string{"invoice_number"})
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "signing root missing")
}

func TestService_CreateProofsInvalidField(t *testing.T) {
	defer setIdentityService(identity.IDService)
	invSrv := getServiceWithMockedLayers()
	i, err := createAnchoredMockDocument(t, false)
	assert.Nil(t, err)
	idService := mockSignatureCheck(i, invSrv)
	setIdentityService(idService)
	_, err = invSrv.CreateProofs(i.CoreDocument.DocumentIdentifier, []string{"invalid_field"})
	assert.Error(t, err)
	assert.Equal(t, "createProofs error No such field: invalid_field in obj", err.Error())
}

func TestService_CreateProofsDocumentDoesntExist(t *testing.T) {
	invService := getServiceWithMockedLayers()
	_, err := invService.CreateProofs(utils.RandomSlice(32), []string{"invoice_number"})
	assert.Error(t, err)
	assert.Equal(t, "document not found: leveldb: not found", err.Error())
}

func TestService_CreateProofsForVersion(t *testing.T) {
	defer setIdentityService(identity.IDService)
	invSrv := getServiceWithMockedLayers()
	i, err := createAnchoredMockDocument(t, false)
	assert.Nil(t, err)
	idService := mockSignatureCheck(i, invSrv)
	setIdentityService(idService)
	olderVersion := i.CoreDocument.CurrentVersion
	i, err = updatedAnchoredMockDocument(t, i)
	assert.Nil(t, err)
	proof, err := invSrv.CreateProofsForVersion(i.CoreDocument.DocumentIdentifier, olderVersion, []string{"invoice_number"})
	assert.Nil(t, err)
	assert.Equal(t, i.CoreDocument.DocumentIdentifier, proof.DocumentId)
	assert.Equal(t, olderVersion, proof.VersionId)
	assert.Equal(t, len(proof.FieldProofs), 1)
	assert.Equal(t, proof.FieldProofs[0].GetProperty(), "invoice_number")
}

func TestService_CreateProofsForVersionDocumentDoesntExist(t *testing.T) {
	i, err := createAnchoredMockDocument(t, false)
	invSrv := getServiceWithMockedLayers()
	assert.Nil(t, err)
	_, err = invSrv.CreateProofsForVersion(i.CoreDocument.DocumentIdentifier, utils.RandomSlice(32), []string{"invoice_number"})
	assert.Error(t, err)
	assert.Equal(t, "document not found for the given version: leveldb: not found", err.Error())
}

func TestService_RequestDocumentSignature_SigningRootNil(t *testing.T) {
	defer setIdentityService(identity.IDService)
	invSrv := getServiceWithMockedLayers()
	i, err := createAnchoredMockDocument(t, true)
	assert.Nil(t, err)
	idService := mockSignatureCheck(i, invSrv)
	setIdentityService(idService)
	i.CoreDocument.SigningRoot = nil
	signature, err := invSrv.RequestDocumentSignature(i)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), strconv.Itoa(int(code.DocumentInvalid)))
	assert.Contains(t, err.Error(), "signing root missing")
	assert.Nil(t, signature)
}

func createAnchoredMockDocument(t *testing.T, skipSave bool) (*Invoice, error) {
	i := &Invoice{
		InvoiceNumber: "test_invoice",
		GrossAmount:   60,
		CoreDocument:  coredocument.New(),
	}
	err := i.calculateDataRoot()
	if err != nil {
		return nil, err
	}
	// get the coreDoc for the invoice
	corDoc, err := i.PackCoreDocument()
	if err != nil {
		return nil, err
	}
	assert.Nil(t, coredocument.FillSalts(corDoc))
	err = coredocument.CalculateSigningRoot(corDoc)
	if err != nil {
		return nil, err
	}

	sig := signatures.Sign(&config.IdentityConfig{
		ID:         centID,
		PublicKey:  key1Pub[:],
		PrivateKey: key1,
	}, corDoc.SigningRoot)

	corDoc.Signatures = append(corDoc.Signatures, sig)

	err = coredocument.CalculateDocumentRoot(corDoc)
	if err != nil {
		return nil, err
	}
	err = i.UnpackCoreDocument(corDoc)
	if err != nil {
		return nil, err
	}

	if !skipSave {
		err = getRepository().Create(i.CoreDocument.CurrentVersion, i)
		if err != nil {
			return nil, err
		}
	}

	return i, nil
}

func updatedAnchoredMockDocument(t *testing.T, i *Invoice) (*Invoice, error) {
	i.GrossAmount = 50
	err := i.calculateDataRoot()
	if err != nil {
		return nil, err
	}
	// get the coreDoc for the invoice
	corDoc, err := i.PackCoreDocument()
	if err != nil {
		return nil, err
	}
	// hacky update to version
	corDoc.CurrentVersion = corDoc.NextVersion
	corDoc.NextVersion = utils.RandomSlice(32)
	if err != nil {
		return nil, err
	}
	err = coredocument.CalculateSigningRoot(corDoc)
	if err != nil {
		return nil, err
	}
	err = coredocument.CalculateDocumentRoot(corDoc)
	if err != nil {
		return nil, err
	}
	err = i.UnpackCoreDocument(corDoc)
	if err != nil {
		return nil, err
	}
	err = getRepository().Create(i.CoreDocument.CurrentVersion, i)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func TestService_DeriveFromUpdatePayload(t *testing.T) {
	invSrv := getServiceWithMockedLayers()
	// nil payload
	doc, err := invSrv.DeriveFromUpdatePayload(nil, nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid payload")
	assert.Nil(t, doc)

	// messed up identifier
	contextHeader, err := documents.NewContextHeader()
	assert.Nil(t, err)
	payload := &clientinvoicepb.InvoiceUpdatePayload{Identifier: "some identifier"}
	doc, err = invSrv.DeriveFromUpdatePayload(payload, contextHeader)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode identifier")
	assert.Nil(t, doc)

	// missing last version
	id := utils.RandomSlice(32)
	payload.Identifier = hexutil.Encode(id)
	doc, err = invSrv.DeriveFromUpdatePayload(payload, contextHeader)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to fetch old version")
	assert.Nil(t, doc)

	// failed to load from data
	old := new(Invoice)
	err = old.InitInvoiceInput(testingdocuments.CreateInvoicePayload(), contextHeader)
	assert.Nil(t, err)
	old.CoreDocument.DocumentIdentifier = id
	old.CoreDocument.CurrentVersion = id
	old.CoreDocument.DocumentRoot = utils.RandomSlice(32)
	err = getRepository().Create(id, old)
	assert.Nil(t, err)
	payload.Data = &clientinvoicepb.InvoiceData{
		Sender:      "0x010101010101",
		Recipient:   "0x010203040506",
		Payee:       "0x010203020406",
		GrossAmount: 42,
		ExtraData:   "some data",
		Currency:    "EUR",
	}
	doc, err = invSrv.DeriveFromUpdatePayload(payload, contextHeader)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load invoice from data")
	assert.Nil(t, doc)

	// failed core document new version
	payload.Data.ExtraData = hexutil.Encode(utils.RandomSlice(32))
	payload.Collaborators = []string{"some wrong ID"}
	doc, err = invSrv.DeriveFromUpdatePayload(payload, contextHeader)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to prepare new version")
	assert.Nil(t, doc)

	// success
	wantCollab := utils.RandomSlice(6)
	payload.Collaborators = []string{hexutil.Encode(wantCollab)}
	doc, err = invSrv.DeriveFromUpdatePayload(payload, contextHeader)
	assert.Nil(t, err)
	assert.NotNil(t, doc)
	cd, err := doc.PackCoreDocument()
	assert.Nil(t, err)
	assert.Equal(t, wantCollab, cd.Collaborators[1])
	assert.Len(t, cd.Collaborators, 2)
	oldCD, err := old.PackCoreDocument()
	assert.Nil(t, err)
	assert.Equal(t, oldCD.DocumentIdentifier, cd.DocumentIdentifier)
	assert.Equal(t, payload.Identifier, hexutil.Encode(cd.DocumentIdentifier))
	assert.Equal(t, oldCD.CurrentVersion, cd.PreviousVersion)
	assert.Equal(t, oldCD.NextVersion, cd.CurrentVersion)
	assert.NotNil(t, cd.NextVersion)
	assert.Equal(t, payload.Data, doc.(*Invoice).getClientData())
}

func TestService_Update(t *testing.T) {
	invSrv := service{repo: getRepository()}
	ctx := context.Background()
	ctxh, err := documents.NewContextHeader()
	assert.Nil(t, err)

	// pack failed
	model := &mockModel{}
	model.On("PackCoreDocument").Return(nil, fmt.Errorf("pack error")).Once()
	_, err = invSrv.Update(ctx, model)
	model.AssertExpectations(t)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "pack error")

	// missing last version
	model = &mockModel{}
	cd := coredocument.New()
	model.On("PackCoreDocument").Return(cd, nil).Once()
	_, err = invSrv.Update(ctx, model)
	model.AssertExpectations(t)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "document not found")

	payload := testingdocuments.CreateInvoicePayload()
	payload.Collaborators = []string{"0x010203040506"}
	inv, err := invSrv.DeriveFromCreatePayload(payload, ctxh)
	assert.Nil(t, err)
	cd, err = inv.PackCoreDocument()
	assert.Nil(t, err)
	cd.DocumentRoot = utils.RandomSlice(32)
	inv.(*Invoice).CoreDocument = cd
	getRepository().Create(cd.CurrentVersion, inv)

	// calculate data root fails
	model = &mockModel{}
	model.On("PackCoreDocument").Return(cd, nil).Once()
	_, err = invSrv.Update(ctx, model)
	model.AssertExpectations(t)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown document type")

	// anchor fails
	data, err := invSrv.DeriveInvoiceData(inv)
	assert.Nil(t, err)
	data.GrossAmount = 100
	data.ExtraData = hexutil.Encode(utils.RandomSlice(32))
	collab := hexutil.Encode(utils.RandomSlice(6))
	newInv, err := invSrv.DeriveFromUpdatePayload(&clientinvoicepb.InvoiceUpdatePayload{
		Identifier:    hexutil.Encode(cd.DocumentIdentifier),
		Collaborators: []string{collab},
		Data:          data,
	}, ctxh)
	assert.Nil(t, err)
	newData, err := invSrv.DeriveInvoiceData(newInv)
	assert.Nil(t, err)
	assert.Equal(t, data, newData)
	proc := &testingutils.MockCoreDocumentProcessor{}
	proc.On("PrepareForSignatureRequests", newInv).Return(nil).Once()
	proc.On("RequestSignatures", ctx, newInv).Return(nil).Once()
	proc.On("PrepareForAnchoring", newInv).Return(nil).Once()
	proc.On("AnchorDocument", newInv).Return(nil).Once()
	proc.On("SendDocument", ctx, newInv).Return(nil).Once()
	invSrv.coreDocProcessor = proc
	inv, err = invSrv.Update(ctx, newInv)
	proc.AssertExpectations(t)
	assert.Nil(t, err)
	assert.NotNil(t, inv)

	newCD, err := inv.PackCoreDocument()
	assert.Nil(t, err)
	assert.True(t, getRepository().Exists(newCD.DocumentIdentifier))
	assert.True(t, getRepository().Exists(newCD.CurrentVersion))
	assert.True(t, getRepository().Exists(newCD.PreviousVersion))

	newData, err = invSrv.DeriveInvoiceData(inv)
	assert.Nil(t, err)
	assert.Equal(t, data, newData)

}

func TestService_calculateDataRoot(t *testing.T) {
	invSrv := getServiceWithMockedLayers().(service)
	ctxh, err := documents.NewContextHeader()
	assert.Nil(t, err)

	// type mismatch
	inv, err := invSrv.calculateDataRoot(nil, &testingdocuments.MockModel{}, nil)
	assert.Nil(t, inv)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unknown document type")

	// failed validator
	inv, err = invSrv.DeriveFromCreatePayload(testingdocuments.CreateInvoicePayload(), ctxh)
	assert.Nil(t, err)
	assert.Nil(t, inv.(*Invoice).CoreDocument.DataRoot)
	v := documents.ValidatorFunc(func(_, _ documents.Model) error {
		return fmt.Errorf("validations fail")
	})
	inv, err = invSrv.calculateDataRoot(nil, inv, v)
	assert.Nil(t, inv)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validations fail")

	// create failed
	inv, err = invSrv.DeriveFromCreatePayload(testingdocuments.CreateInvoicePayload(), ctxh)
	assert.Nil(t, err)
	assert.Nil(t, inv.(*Invoice).CoreDocument.DataRoot)
	err = invSrv.repo.Create(inv.(*Invoice).CoreDocument.CurrentVersion, inv)
	assert.Nil(t, err)
	inv, err = invSrv.calculateDataRoot(nil, inv, CreateValidator())
	assert.Nil(t, inv)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "document already exists")

	// success
	inv, err = invSrv.DeriveFromCreatePayload(testingdocuments.CreateInvoicePayload(), ctxh)
	assert.Nil(t, err)
	assert.Nil(t, inv.(*Invoice).CoreDocument.DataRoot)
	inv, err = invSrv.calculateDataRoot(nil, inv, CreateValidator())
	assert.Nil(t, err)
	assert.NotNil(t, inv)
	assert.NotNil(t, inv.(*Invoice).CoreDocument.DataRoot)
}
