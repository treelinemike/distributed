package gfcommon

type (
	NetworkAddress struct {
		Address string `mapstructure:"Address"`
		Port    string `mapstructure:"Port"`
	}
	GFPlayerConfig struct {
		Host         NetworkAddress
		OtherPlayers []NetworkAddress
	}
)
