---
Profiles:
    GenesisProfile:
        Orderer:
            <<: *OrdererDefaults
            Organizations:
                - *OrdererOrg
        Consortiums:
            SampleConsortium:
                Organizations:
                    {{range $peerorg := .PeerOrgs}}- *{{$peerorg.NameHash}}
                    {{end}}
    
    {{range $channel := .Channels}}
    {{$channel.Name}}Profile:
        Consortium: SampleConsortium
        Application:
            <<: *ApplicationDefaults
            Organizations:
                {{range $channelorgs := $channel.OrgNameHash}}- *{{$channelorgs}}
                {{end}}
    {{end}}

Organizations:
    - &OrdererOrg
        Name: {{.OrdererOrg.NameHash}}
        ID: {{.OrdererOrg.NameHash}}MSP
        MSPDir: crypto-config/ordererOrganizations/{{$.OrdererOrg.NameHash}}.orderer-{{$.OrdererOrg.NameHash}}.{{$.Domain.Namespace}}.{{$.Domain.Surfix}}/msp

{{range $peerorg := .PeerOrgs}}
    - &{{$peerorg.NameHash}}
        Name: {{$peerorg.NameHash}}
        ID: {{$peerorg.NameHash}}MSP
        MSPDir: crypto-config/peerOrganizations/{{$peerorg.NameHash}}.peer-{{$peerorg.NameHash}}.{{$.Domain.Namespace}}.{{$.Domain.Surfix}}/msp
        AnchorPeers:
            - Host: peer-{{$peerorg.NameHash}}-0.peer-{{$peerorg.NameHash}}.{{$.Domain.Namespace}}.{{$.Domain.Surfix}}
              Port: 7051
{{end}}

Orderer: &OrdererDefaults
    OrdererType: {{ if eq .Consensus "FBFT" }}fbft{{else if eq .Consensus "PBFT" }}pbft{{else if eq .Consensus "SFLIC" }}sflic{{else if eq .Consensus "MFLIC" }}mflic{else if eq .Consensus "kafka" }}kafka{{else}}{{.Consensus}}{{end}}
    Addresses:
        {{range $index := .OrdererOrg.NodeNumList}}- orderer-{{$.OrdererOrg.NameHash}}-{{$index}}.orderer-{{$.OrdererOrg.NameHash}}.{{$.Domain.Namespace}}.{{$.Domain.Surfix}}:7050
        {{end}}
    BatchTimeout: 2s
    BatchSize:
        MaxMessageCount: 10
        AbsoluteMaxBytes: 99 MB
        PreferredMaxBytes: 512 KB
    Kafka:
        # Brokers: A list of Kafka brokers to which the orderer connects
        # NOTE: Use IP:port notation
        Brokers:
            {{range $index := .KafkaOrg.NodeNumList}}- kafka-{{$.KafkaOrg.NameHash}}-{{$index}}.kafka-{{$.KafkaOrg.NameHash}}.{{$.Domain.Namespace}}.{{$.Domain.Surfix}}:9093
            {{end}}
    Organizations:

Application: &ApplicationDefaults
    Organizations:
