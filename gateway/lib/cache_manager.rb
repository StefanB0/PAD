require 'singleton'
require 'digest'
require 'logger'
require 'redis'
require 'timeout'

class CacheManager
  include Singleton
  
  def initialize
    @nodes_list = {}
    @server_list = {}
    @cache = {}

    @logger = Logger.new(STDERR)
  end

  def add_node(node, address)
    redis = Redis.new(host: address.split(':')[0], port: address.split(':')[1], timeout: 1)
    
    @server_list[node] = redis

    labels = create_labels(node)
    hashes = labels.map { |label| hash_value(label) }
    hashes.each { |hash| @nodes_list[hash] = node }

    remap_caches
  end

  def remove_node(node)
    @logger.info("Removing node #{node}")

    labels = create_labels(node)
    hashes = labels.map { |label| hash_value(label) }
    hashes.each { |hash| @nodes_list.delete(hash) }

    if check_node_status(node)
      @logger.info("Remapping node #{node}")
      @logger.info("Node #{node} has #{count_node_keys(node)} keys")
      remap_caches
    else
      @cache.each do |key, value|
        if value == node
          @cache.delete(key)
        end
      end
    end

    @server_list.delete(node)
  end

  def remap_caches
    remaps = 0
    @cache.each do |key, server|
      node = find_cache(key)
      if node != server
        @cache[key] = node
        value = @server_list[server].get(key)
        @server_list[node].set(key, value)
        remaps += 1
      end
    end
    @logger.info("Remapped #{remaps} keys")
    remaps
  end

  def check_node_status(node)
    if @server_list[node].nil?
      @logger.error("Node #{node} already removed")
      return false
    end

    status = ping_with_timeout(@server_list[node], 1)
    if status
      return true
    else
      @logger.error("Node #{node} offline")
      return false
    end

    begin
      Timeout.timeout(1) do
        @logger.info("Checking node #{node} status")
        @server_list[node].ping
        return true
      end
    rescue Exception => e
      @logger.error("Node #{node} offline. Error: #{e}")
      false
    end
  end

  def delete_all_cache
    @server_list.each do |node, server|
      check_node_status(node) ? server.flushdb : next
      server.flushdb
    end
    @cache = {}
  end
  
  def get_server(node)
    @server_list[node]
  end

  def get_nodes
    @nodes_list.clone
  end

  def get_cache
    @cache.clone
  end

  def create_labels(node)
    (1..10).map { |i| "#{node}:#{i}" }
  end

  def hash_value(node)
    Digest::SHA256.digest(node).sum % 360
  end

  def add_key(key, value=42)
    node = find_cache(key)
    @logger.info("Key #{key} hash is #{hash_value(key)}")
    @logger.info("Adding key #{key} to node #{node}")
    
    if check_node_status(node)
      @cache[key] = node
      @server_list[node].set(key, value)
      @server_list[node].expire(key, 300) # Redis timeout 5 minutes
    else
      remove_node(node)
      add_key(key, value)
    end
  end

  def get_key(key)
    node = get_key_node(key)

    if node.nil?
      nil
    elsif check_node_status(node)
      @logger.info "Key #{key} retrieved from node #{node}"
      @server_list[node].get(key)
    else
      Thread.new { remove_node(node) }
      # remove_node(node)
      nil
    end
  end

  def remove_key(key)
    node = get_key_node(key)

    if check_node_status(node)
      @server_list[node].del(key)
      @cache.delete(key)
    else
      remove_node(node)
      remove_key(key)
    end
  end

  def get_key_node(key)
    @cache[key]
  end

  def count_node_keys(node)
    @cache.values.count(node)
  end

  def find_cache(key)
    hash = hash_value(key)

    node_hash = closest_node_hash(hash)
    @nodes_list[node_hash]
  end

  def closest_node_hash(key)
    @nodes_list.keys.sort.bsearch { |server| server >= key } || @nodes_list.keys.sort.first
  end

  def ping_with_timeout(redis, timeout)
    thread = Thread.new do
      begin
        redis.ping
      rescue Redis::BaseError => e
        "Error: #{e.message}"
      end
    end
  
    thread.join(timeout)
  
    if thread.alive?
      thread.kill
      return false
    end
  
    true
  end
end

cache_manager = CacheManager.instance

cache_manager.add_node("Redis1", ENV['REDIS1_ADDRESS'])
cache_manager.add_node("Redis2", ENV['REDIS2_ADDRESS'])
cache_manager.add_node("Redis3", ENV['REDIS3_ADDRESS'])

cache_manager.delete_all_cache

# keys = ["Alex", "Bob", "Charlie", "David", "Ethan", "Frank", "George"]

# keys.each { |key| cache_manager.add_key(key, rand(1..100)) }
# # cache_manager.remove_node(s3)
# cache_manager.remove_node(s4)

# keys.each { |key| puts "#{key} => #{cache_manager.get_key(key)}" }

