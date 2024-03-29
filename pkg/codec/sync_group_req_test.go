package codec

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDecodeIllegalSyncGroupReq(t *testing.T) {
	bytes := make([]byte, 0)
	_, err := DecodeSyncGroupReq(bytes, 0)
	assert.NotNil(t, err)
}

func TestDecodeSyncGroupReqV4(t *testing.T) {
	bytes := testHex2Bytes(t, "00000006002f636f6e73756d65722d38646437623936622d366239342d346139622d623263632d3363623538393863396364662d31002538646437623936622d366239342d346139622d623263632d3363623538393863396364660000000155636f6e73756d65722d38646437623936622d366239342d346139622d623263632d3363623538393863396364662d312d34333361636236612d653665632d343561612d623738642d366132343963666630376663000255636f6e73756d65722d38646437623936622d366239342d346139622d623263632d3363623538393863396364662d312d34333361636236612d653665632d343561612d623738642d3661323439636666303766631b0001000000010006746573742d350000000100000000ffffffff0000")
	fetchReq, err := DecodeSyncGroupReq(bytes, 4)
	assert.Nil(t, err)
	assert.Equal(t, 6, fetchReq.CorrelationId)
	assert.Equal(t, "consumer-8dd7b96b-6b94-4a9b-b2cc-3cb5898c9cdf-1", fetchReq.ClientId)
	assert.Equal(t, "8dd7b96b-6b94-4a9b-b2cc-3cb5898c9cdf", fetchReq.GroupId)
	assert.Equal(t, 1, fetchReq.GenerationId)
	assert.Equal(t, "8dd7b96b-6b94-4a9b-b2cc-3cb5898c9cdf", fetchReq.GroupId)
	assert.Equal(t, "", fetchReq.ProtocolType)
	assert.Equal(t, "", fetchReq.ProtocolName)
	assert.Len(t, fetchReq.GroupAssignments, 1)
	groupAssignment := fetchReq.GroupAssignments[0]
	assert.Equal(t, "consumer-8dd7b96b-6b94-4a9b-b2cc-3cb5898c9cdf-1-433acb6a-e6ec-45aa-b78d-6a249cff07fc", groupAssignment.MemberId)
}
