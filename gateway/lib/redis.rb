require 'redis'
require 'singleton'
require 'dotenv'

Dotenv.load('.env')

class RedisCache 
  include Singleton

  def initialize(host="localhost", port=6379)
    @redis = Redis.new(:timeout => ENV['REDIS_TIMEOUT'].to_i)
    redis = Redis.new(host: host, port: port, timeout: ENV['REDIS_TIMEOUT'].to_i)
  end

  def check_token(access_token)
    @redis.get(access_token).nil? ? false : true
  end

  def set(key, value)
    @redis.set(key, value)
  end

  def get(key)
    @redis.get(key)
  end

  def delete(key)
    @redis.del(key)
  end

  def flush
    @redis.flushall
  end
end

