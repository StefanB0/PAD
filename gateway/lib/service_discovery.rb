require 'singleton'
require 'net/http'
require 'dotenv'

Dotenv.load('.env')

class ServiceDiscovery
  include Singleton

  attr_accessor :sd_address, :sd_port

  def initialize
    @sd_address = ENV['SD_ADDRESS']
    @sd_port = ENV['SD_PORT']
  end

  def get_service_address(service_name)
    uri = URI("http://#{@sd_address}:#{@sd_port}/getservice/#{service_name}")
    Net::HTTP.get(uri)
  end
end
