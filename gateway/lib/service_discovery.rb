require 'singleton'
require 'net/http'
require 'dotenv'
require 'json'

Dotenv.load('.env')

class ServiceDiscovery
  include Singleton

  def initialize
    @sd_address = ENV['SD_ADDRESS']
  end

  def get_service_address(service_name)
    uri = URI("http://#{@sd_address}/service/#{service_name}")
    Net::HTTP.get(uri)
  end

  def subscribe 
    uri = URI("http://#{@sd_address}/service")
    # response = Net::HTTP.post_form(uri, 'name' => 'gateway_service', 'address' => 'http://localhost:8080')
    http = Net::HTTP.new(uri.host, uri.port)
    request = Net::HTTP::Post.new(uri, initheader = {'Content-Type' =>'application/json'})
    request.body = {name: 'gateway_service', address: 'http://localhost:8080'}.to_json
    response = http.request(request)

    @secretkey = JSON.parse(response.body)['secret_key']
  end

  def unsubscribe
    uri = URI("http://#{@sd_address}/service/gateway_service")
    Net::HTTP.start(uri.host, uri.port) do |http|
      request = Net::HTTP::Delete.new(uri)
      request.set_form_data({'address' => 'http://localhost:8080', 'secretkey' => @secretkey})
      http.request(request)
    end    
  end
end
