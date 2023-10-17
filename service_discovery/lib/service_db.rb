require 'singleton'
require 'securerandom'

class ServiceDB 
  include Singleton
  
  def initialize
    @db = {}
    @secretdb = {}
  end
  
  def get_service(name)
    @db[name].nil? ? :invalid_service : @db[name]
  end
  
  def add_service(name, address)
    secretkey = SecureRandom.hex(16)
    if @db[name].nil?
      @db[name] = [address] 
    else
      @db[name].push(address).uniq!
    end
    @secretdb[name+address] = secretkey
    secretkey
  end

  def remove_service(name, address, secretkey)
    if name.nil? || address.nil? || secretkey.nil? || name.empty? || address.empty? || secretkey.empty?
      return :unauthorized
    end

    unless @secretdb[name+address] == secretkey
      return :unauthorized
    end

    @db[name].delete(address)
    if @db[name].empty?
      @db.delete(name)
    end
  end
end