require 'redis'

module RedisCache 
  @redis = Redis.new(:timeout => ENV['REDIS_TIMEOUT'])

  def self.setAccessToken(key, value)
    @redis.set(key, value)
    @redis.expire(key, ENV['ACCESS_TOKEN_DURATION'])
  end

  def self.setRefreshToken(key, value)
    @redis.set(key, value)
    @redis.expire(key, ENV['REFRESH_TOKEN_DURATION'])
  end

  def self.get(key)
    @redis.get(key)
  end

  def self.delete(key)
    @redis.del(key)
  end

  def self.flush
    @redis.flushall
  end
end

