package network

import (
	"github.com/paashzj/kafka_go/pkg/codec"
	"github.com/paashzj/kafka_go/pkg/log"
	"github.com/paashzj/kafka_go/pkg/network/context"
	"github.com/paashzj/kafka_go/pkg/service"
	"github.com/panjf2000/gnet"
	"k8s.io/klog/v2"
)

func (s *Server) Fetch(ctx *context.NetworkContext, frame []byte, version int16) ([]byte, gnet.Action) {
	if version == 11 {
		return s.ReactFetchVersion(ctx, frame, version)
	}
	klog.Error("unknown fetch version ", version)
	return nil, gnet.Close
}

func (s *Server) ReactFetchVersion(ctx *context.NetworkContext, frame []byte, version int16) ([]byte, gnet.Action) {
	req, err := codec.DecodeFetchReq(frame, version)
	if err != nil {
		return nil, gnet.Close
	}
	if !s.checkSasl(ctx) {
		return nil, gnet.Close
	}
	log.Codec().Info("fetch req ", req)
	lowReq := &service.FetchReq{}
	lowReq.FetchTopicReqList = make([]*service.FetchTopicReq, len(req.FetchTopics))
	for i, topicReq := range req.FetchTopics {
		if !s.checkSaslTopic(ctx, topicReq.Topic) {
			return nil, gnet.Close
		}
		lowTopicReq := &service.FetchTopicReq{}
		lowTopicReq.Topic = topicReq.Topic
		lowTopicReq.FetchPartitionReqList = make([]*service.FetchPartitionReq, len(topicReq.FetchPartitions))
		for j, partitionReq := range topicReq.FetchPartitions {
			lowPartitionReq := &service.FetchPartitionReq{}
			lowPartitionReq.PartitionId = partitionReq.PartitionId
			lowPartitionReq.FetchOffset = partitionReq.FetchOffset
			lowTopicReq.FetchPartitionReqList[j] = lowPartitionReq
		}
		lowReq.FetchTopicReqList[i] = lowTopicReq
	}
	lowTopicRespList, err := service.Fetch(ctx.Addr, s.kafkaImpl, lowReq)
	if err != nil {
		return nil, gnet.Close
	}
	resp := codec.NewFetchResp(req.CorrelationId)
	resp.TopicResponses = make([]*codec.FetchTopicResp, len(lowTopicRespList))
	for i, lowTopicResp := range lowTopicRespList {
		f := &codec.FetchTopicResp{}
		f.Topic = lowTopicResp.Topic
		f.PartitionDataList = make([]*codec.FetchPartitionResp, len(lowTopicResp.FetchPartitionRespList))
		for j, p := range lowTopicResp.FetchPartitionRespList {
			partitionResp := &codec.FetchPartitionResp{}
			partitionResp.PartitionIndex = p.PartitionId
			partitionResp.ErrorCode = 0
			partitionResp.HighWatermark = p.HighWatermark
			partitionResp.LastStableOffset = p.LastStableOffset
			partitionResp.LogStartOffset = p.LogStartOffset
			partitionResp.RecordBatch = &codec.RecordBatch{}
			s.convertRecordBatch(partitionResp.RecordBatch, p.RecordBatch)
			partitionResp.AbortedTransactions = -1
			partitionResp.ReplicaData = -1
			f.PartitionDataList[j] = partitionResp
		}
		resp.TopicResponses[i] = f
	}
	return resp.Bytes(version), gnet.None
}

func (s *Server) convertRecordBatch(recordBatch *codec.RecordBatch, lowRecordBatch *service.RecordBatch) {
	recordBatch.Offset = lowRecordBatch.Offset
	recordBatch.MessageSize = lowRecordBatch.MessageSize
	recordBatch.LeaderEpoch = 0
	recordBatch.MagicByte = 2
	recordBatch.Flags = 0
	recordBatch.LastOffsetDelta = lowRecordBatch.LastOffsetDelta
	recordBatch.FirstTimestamp = lowRecordBatch.FirstTimestamp
	recordBatch.LastTimestamp = lowRecordBatch.LastTimestamp
	recordBatch.ProducerId = -1
	recordBatch.ProducerEpoch = -1
	recordBatch.BaseSequence = lowRecordBatch.BaseSequence
	recordBatch.Records = make([]*codec.Record, len(lowRecordBatch.Records))
	for i, r := range lowRecordBatch.Records {
		record := &codec.Record{}
		record.RecordAttributes = 0
		record.RelativeTimestamp = r.RelativeTimestamp
		record.RelativeOffset = r.RelativeOffset
		record.Key = r.Key
		record.Value = r.Value
		record.Headers = r.Headers
		recordBatch.Records[i] = record
	}
}
