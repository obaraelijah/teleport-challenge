package config

// Note: generally I would avoid having a "config.go" as a place for a bunch of
//
//	unrelated constants.  However, here these constants represent the
//	hard-coded values for the exercise, so I thought this might make it
//	easier to adjust the values for experimentation if I put them all in
//	one place.
const (
	CgexecPath = "/tmp/cgexec"

	CgroupDefaultCpuLimit        = 0.5
	CgroupDefaultMemoryLimit     = "2M"
	CgroupDefaultBlkioDevice     = "8:16"
	CgroupDefaultBlkioWriteLimit = CgroupDefaultBlkioDevice + " 20971520"
	CgroupDefaultBlkioReadLimit  = CgroupDefaultBlkioDevice + " 41943040"
)
