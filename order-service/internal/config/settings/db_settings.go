package settings

type DbSettings struct {
	ConnectionString          string `mapstructure:"ConnectionString"`
	MigrationConnectionString string `mapstructure:"MigrationConnectionString"`
}
