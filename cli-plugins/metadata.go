package cliplugins

var (
	// NamePrefix is the prefix required on all plugin binary names
	NamePrefix = "docker-"

	// MetadataSubcommandName is the name of the plugin subcommand
	// which must be supported by every plugin and returns the
	// plugin metadata.
	MetadataSubcommandName = "docker-cli-plugin-metadata"
)

// Metadata provided by the plugin
type Metadata struct {
	Version          string
	Vendor           string
	ShortDescription string
	SchemaVersion    string
}
