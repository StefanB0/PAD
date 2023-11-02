require 'singleton'
require 'securerandom'
require 'httparty'

class ServiceDB 
  include Singleton
  
  def initialize
    @db = {}
    @secretdb = {}
  end
  
  def get_service(name)
    @db[name].nil? ? :invalid_service : @db[name]
  end

  def get_all_services
    services = []
    @db.each do |name, addresses|
      services.push({name: name, addresses: addresses})
    end
    services
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

  def flush_offline_services
    @db.each do |name, addresses|
      addresses.each do |address|
        unless ping_service(address)
          puts "Service #{name} at #{address} is offline. Removing from database."
          @db[name].delete(address)
          @secretdb.delete(name+address)
        end
      end
    end
  end

  def ping_service(address)
    success = false
    begin
      response = HTTParty.get(address + '/status')
      success = response.code == 200 && response.body == 'OK'
    rescue Exception => e
      puts e
    end
    return success
  end
end