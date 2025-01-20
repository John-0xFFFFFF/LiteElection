package election

type keepAliveResult int64
const(
	keyNotFound keepAliveResult=iota
	keepAliveSucceeded
	leaderChanged
)

var keepAliveScript=`
local key = KEYS[1]
local inputValue = ARGV[1]
local expireTime = tonumber(ARGV[2])

local redisValue = redis.call("GET", key)

if redisValue == inputValue then
    redis.call("PEXPIRE", key, expireTime)
    return 1
elseif redisValue == false then
    return 0
else
    return 2
end
`

var setNxExScript = `
local current_value = redis.call('GET', KEYS[1])
if current_value then
	return current_value
else
	redis.call('SET', KEYS[1], ARGV[1], 'EX', ARGV[2], 'NX')
	return "OK"
end
`