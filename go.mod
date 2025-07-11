module github.com/shirou/gopsutil/v3

go 1.18

require (
	github.com/google/go-cmp v0.5.9
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0
	github.com/pkg/errors v0.9.1
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c
	github.com/stretchr/testify v1.8.2
	github.com/tklauser/go-sysconf v0.3.11
	github.com/yusufpapurcu/wmi v1.2.2
	go.easyops.local/slog v0.0.0-20240218012623-b8d04ff912aa
	golang.org/x/sys v0.6.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/tklauser/numcpus v0.6.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract v3.22.11
