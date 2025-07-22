package config

//定义一个全局变量

var Config *ConfigStruct

type ConfigStruct struct {
	//变量名  +  类型
	Mysql MysqlConfig
}

//定义所有的config的结构体

// mysql的配置
type MysqlConfig struct {
	// 变量名  +  类型  + 标签  toml是配置文件的标签
	//变量名的首字母大写意味着可以被外部引用
	Address      string `toml:"address"`
	Port         string `toml:"port"`
	DbName       string `toml:"db_name"`
	UserName     string `toml:"user_name"`
	Password     string `toml:"password"`
	MaxOpenConns int    `toml:"max_open_conns"`
	MaxIdleConns int    `toml:"max_idle_conns"`
	MaxLifeTime  int    `toml:"max_life_time"`
}

// 用于存储测试网（TestNet）区块链的相关配置信息。
// ChainId：链的 ID（区块链网络的唯一标识）
// NetUrl：区块链节点的访问地址（RPC URL）
// PlgrAddress：PLGR 代币的合约地址
// PledgePoolToken：质押池合约的 Token 地址
// BscPledgeOracleToken：BSC 质押预言机合约的 Token 地址
type TestNetConfig struct {
	ChainId              string `toml:"chain_id"`
	NetUrl               string `toml:"net_url"`
	PlgrAddress          string `toml:"plgr_address"`
	PledgePoolToken      string `toml:"pledge_pool_token"`
	BscPledgeOracleToken string `toml:"bsc_pledge_oracle_token"`
}

// 用于存储主网（MainNet）区块链的相关配置信息。
type MainNetConfig struct {
	ChainId              string `toml:"chain_id"`
	NetUrl               string `toml:"net_url"`
	PlgrAddress          string `toml:"plgr_address"`
	PledgePoolToken      string `toml:"pledge_pool_token"`
	BscPledgeOracleToken string `toml:"bsc_pledge_oracle_token"`
}

type RedisConfig struct {
	Address     string `toml:"address"`
	Port        string `toml:"port"`
	Db          int    `toml:"db"`
	Password    string `toml:"password"`
	MaxIdle     int    `toml:"max_idle"`
	MaxActive   int    `toml:"max_active"`
	IdleTimeout int    `toml:"idle_timeout"`
}
