package config

// Domain 站点主域名
const Domain = ""

// FetchCdnFromGithub 从Github上获取cdn列表
const FetchCdnFromGithub = true

// CdnFetchURI 获取cdn的地址
const CdnFetchURI = "https://api.github.com/repos/testerSunshine/12306/contents/cdn_list"

// Verify12306URI 12306验证用的地址
const Verify12306URI = "otn/zwdch/init"

// Server12306URL 12306客运服务地址
const Server12306URL = "kyfw.12306.cn"

// CdnDealCountPerTime 每次处理校验cdn的数量
const CdnDealCountPerTime = 100
