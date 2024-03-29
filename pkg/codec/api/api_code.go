package api

type Code int16

const (
	Produce Code = iota
	Fetch
	ListOffsets
	Metadata
	LeaderAndIsr
	StopReplica
	UpdateMetadata
	ControlledShutdown
	OffsetCommit
	OffsetFetch
	FindCoordinator
	JoinGroup
	Heartbeat
	LeaveGroup
	SyncGroup
	DescribeGroups
	ListGroups
	SaslHandshake
	ApiVersions
	CreateTopics
	DeleteTopics
	DeleteRecords
	InitProducerId
	OffsetForLeaderEpoch
	AddPartitionsToTxn
	AddOffsetsToTxn
	EndTxn
	WriteTxnMarkers
	TxnOffsetCommit
	DescribeAcls
	CreateAcls
	DeleteAcls
	DescribeConfigs
	AlterConfigs
	AlterReplicaLogDirs
	DescribeLogDirs
	SaslAuthenticate
	CreatePartitions
	CreateDelegationToken
	RenewDelegationToken
	ExpireDelegationToken
	DescribeDelegationToken
	DeleteGroups
	ElectLeaders
	IncrementalAlterConfigs
	AlterPartitionReassignments
	ListPartitionReassignments
	OffsetDelete
	DescribeClientQuotas
	AlterClientQuotas
	DescribeUserScramCredentials
	AlterUserScramCredentials
	AlterIsr
	UpdateFeatures
	DescribeCluster
	DescribeProducers
)
