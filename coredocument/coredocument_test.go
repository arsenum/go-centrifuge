// +build unit

package coredocument

import (
	"crypto/sha256"
	"testing"

	"github.com/centrifuge/centrifuge-protobufs/documenttypes"
	"github.com/centrifuge/centrifuge-protobufs/gen/go/coredocument"
	"github.com/centrifuge/go-centrifuge/identity"
	"github.com/centrifuge/go-centrifuge/testingutils"
	"github.com/centrifuge/go-centrifuge/utils"
	"github.com/centrifuge/precise-proofs/proofs"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/golang/protobuf/ptypes/any"
	"github.com/stretchr/testify/assert"
)

var (
	id1 = utils.RandomSlice(32)
	id2 = utils.RandomSlice(32)
	id3 = utils.RandomSlice(32)
	id4 = utils.RandomSlice(32)
	id5 = utils.RandomSlice(32)

	centID  = utils.RandomSlice(identity.CentIDLength)
	key1Pub = [...]byte{230, 49, 10, 12, 200, 149, 43, 184, 145, 87, 163, 252, 114, 31, 91, 163, 24, 237, 36, 51, 165, 8, 34, 104, 97, 49, 114, 85, 255, 15, 195, 199}
	key1    = []byte{102, 109, 71, 239, 130, 229, 128, 189, 37, 96, 223, 5, 189, 91, 210, 47, 89, 4, 165, 6, 188, 53, 49, 250, 109, 151, 234, 139, 57, 205, 231, 253, 230, 49, 10, 12, 200, 149, 43, 184, 145, 87, 163, 252, 114, 31, 91, 163, 24, 237, 36, 51, 165, 8, 34, 104, 97, 49, 114, 85, 255, 15, 195, 199}
)

func TestGetSigningProofHashes(t *testing.T) {
	docAny := &any.Any{
		TypeUrl: documenttypes.InvoiceDataTypeUrl,
		Value:   []byte{},
	}
	cd := New()
	cd.EmbeddedData = docAny
	cd.DataRoot = utils.RandomSlice(32)
	cds := &coredocumentpb.CoreDocumentSalts{}
	err := proofs.FillSalts(cd, cds)
	assert.Nil(t, err)

	cd.CoredocumentSalts = cds
	err = CalculateSigningRoot(cd)
	assert.Nil(t, err)

	err = CalculateDocumentRoot(cd)
	assert.Nil(t, err)

	hashes, err := getSigningProofHashes(cd)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(hashes))

	valid, err := proofs.ValidateProofSortedHashes(cd.SigningRoot, hashes, cd.DocumentRoot, sha256.New())
	assert.True(t, valid)
	assert.Nil(t, err)
}

func TestGetDataProofHashes(t *testing.T) {
	docAny := &any.Any{
		TypeUrl: documenttypes.InvoiceDataTypeUrl,
		Value:   []byte{},
	}
	cd := New()
	cd.EmbeddedData = docAny
	cd.DataRoot = utils.RandomSlice(32)
	cds := &coredocumentpb.CoreDocumentSalts{}
	err := proofs.FillSalts(cd, cds)
	assert.Nil(t, err)

	cd.CoredocumentSalts = cds

	err = CalculateSigningRoot(cd)
	assert.Nil(t, err)

	err = CalculateDocumentRoot(cd)
	assert.Nil(t, err)

	hashes, err := getDataProofHashes(cd)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(hashes))

	valid, err := proofs.ValidateProofSortedHashes(cd.DataRoot, hashes, cd.DocumentRoot, sha256.New())
	assert.True(t, valid)
	assert.Nil(t, err)
}

func TestGetDocumentSigningTree(t *testing.T) {
	docAny := &any.Any{
		TypeUrl: documenttypes.InvoiceDataTypeUrl,
		Value:   []byte{},
	}
	cd := New()
	cd.EmbeddedData = docAny
	cds := &coredocumentpb.CoreDocumentSalts{}
	proofs.FillSalts(cd, cds)
	cd.CoredocumentSalts = cds
	tree, err := GetDocumentSigningTree(cd)
	assert.Nil(t, err)
	assert.NotNil(t, tree)

	_, leaf := tree.GetLeafByProperty("document_type")
	assert.NotNil(t, leaf)
}

func TestGetDocumentSigningTree_EmptyEmbeddedData(t *testing.T) {
	cd := New()
	cds := &coredocumentpb.CoreDocumentSalts{}
	proofs.FillSalts(cd, cds)
	cd.CoredocumentSalts = cds
	tree, err := GetDocumentSigningTree(cd)
	assert.NotNil(t, err)
	assert.Nil(t, tree)
}

// TestGetDocumentRootTree tests that the documentroottree is properly calculated
func TestGetDocumentRootTree(t *testing.T) {
	cd := &coredocumentpb.CoreDocument{SigningRoot: []byte{0x72, 0xee, 0xb8, 0x88, 0x92, 0xf7, 0x6, 0x19, 0x82, 0x76, 0xe9, 0xe7, 0xfe, 0xcc, 0x33, 0xa, 0x66, 0x78, 0xd4, 0xa6, 0x5f, 0xf6, 0xa, 0xca, 0x2b, 0xe4, 0x17, 0xa9, 0xf6, 0x15, 0x67, 0xa1}}
	tree, err := GetDocumentRootTree(cd)

	// Manually constructing the two node tree:
	signaturesLengthLeaf := sha256.Sum256(append([]byte("signatures.length0"), make([]byte, 32)...))
	expectedRootHash := sha256.Sum256(append(signaturesLengthLeaf[:], cd.SigningRoot...))
	assert.Nil(t, err)
	assert.Equal(t, expectedRootHash[:], tree.RootHash())
}

func TestGetTypeUrl(t *testing.T) {
	coreDocument := testingutils.GenerateCoreDocument()

	documentType, err := GetTypeURL(coreDocument)
	assert.Nil(t, err, "should not throw an error because coreDocument has a type")
	assert.NotEqual(t, documentType, "", "document type shouldn't be empty")

	_, err = GetTypeURL(nil)
	assert.Error(t, err, "nil should throw an error")

	coreDocument.EmbeddedData.TypeUrl = ""
	_, err = GetTypeURL(nil)
	assert.Error(t, err, "should throw an error because typeUrl is not set")
}

func TestPrepareNewVersion(t *testing.T) {
	var doc coredocumentpb.CoreDocument
	id := utils.RandomSlice(32)
	cv := id
	nv := utils.RandomSlice(32)
	dr := utils.RandomSlice(32)
	doc = coredocumentpb.CoreDocument{}

	// failed new with collaborators
	collabs := []string{"some ID"}
	newDoc, err := PrepareNewVersion(doc, collabs)
	assert.Error(t, err)
	assert.Nil(t, newDoc)

	// missing doc identifier
	collabs = []string{"0x010203040506"}
	newDoc, err = PrepareNewVersion(doc, collabs)
	assert.NotNil(t, err)
	assert.Nil(t, newDoc)

	//missing current version
	doc.DocumentIdentifier = id
	newDoc, err = PrepareNewVersion(doc, collabs)
	assert.NotNil(t, err)
	assert.Nil(t, newDoc)

	doc.CurrentVersion = cv
	newDoc, err = PrepareNewVersion(doc, collabs)
	assert.NotNil(t, err)
	assert.Nil(t, newDoc)

	doc.NextVersion = nv
	newDoc, err = PrepareNewVersion(doc, collabs)
	assert.NotNil(t, err)
	assert.Nil(t, newDoc)

	doc.CurrentVersion = cv
	doc.NextVersion = nv
	doc.DocumentRoot = dr

	newDoc, err = PrepareNewVersion(doc, collabs)
	assert.Nil(t, err)

	// original document hasn't changed
	assert.Equal(t, cv, doc.CurrentVersion)

	// new document has changed
	assert.Equal(t, id, newDoc.DocumentIdentifier)
	assert.Equal(t, cv, newDoc.PreviousVersion)
	assert.Equal(t, nv, newDoc.CurrentVersion)
	assert.Equal(t, dr, newDoc.PreviousRoot)
	assert.Nil(t, newDoc.DocumentRoot)
}

func TestNewWithCollaborators(t *testing.T) {
	// messed up collaborators
	c := []string{"some id"}
	cd, err := NewWithCollaborators(c)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to decode collaborator")
	assert.Nil(t, cd)

	// success
	c1 := utils.RandomSlice(6)
	c2 := utils.RandomSlice(6)
	c = []string{hexutil.Encode(c1), hexutil.Encode(c2)}
	cd, err = NewWithCollaborators(c)
	assert.Nil(t, err)
	assert.NotNil(t, cd)
	assert.NotNil(t, cd.DocumentIdentifier)
	assert.NotNil(t, cd.CurrentVersion)
	assert.NotNil(t, cd.NextVersion)
	assert.NotNil(t, cd.Collaborators)
	assert.NotNil(t, cd.CoredocumentSalts)
	assert.Equal(t, [][]byte{c1, c2}, cd.Collaborators)
}

func TestGetExternalCollaborators(t *testing.T) {
	c1 := utils.RandomSlice(6)
	c2 := utils.RandomSlice(6)
	c := []string{hexutil.Encode(c1), hexutil.Encode(c2)}
	cd, err := NewWithCollaborators(c)
	assert.Equal(t, [][]byte{c1, c2}, cd.Collaborators)
	collaborators, err := GetExternalCollaborators(cd)
	assert.Nil(t, err)
	assert.NotNil(t, collaborators)
	assert.Equal(t, [][]byte{c1, c2}, collaborators)
}
