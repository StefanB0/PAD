require 'sinatra'
require 'httparty'
require 'net/http'
require 'json'
require 'logger'


require_relative '../lib/redis.rb'
require_relative '../lib/service_discovery.rb'
require_relative '../lib/load_balancer.rb'

load_balancer = LoadBalancer.new('analytics_service')
logger = Logger.new(STDERR)

get '/recommend/:tag' do
  tag = params['tag']
  url = URI("#{load_balancer.next_item}/getRecommendations")
  if url.nil?
    return 404
  end

  req = Net::HTTP::Post.new(url, initheader = {'Content-Type' =>'application/json'})
  req.body = {"tag" => tag}.to_json

  res = Net::HTTP.start(url.hostname, url.port) do |http|
    http.request(req)
  end

  return res.code.to_i if res.code != '200'

  content_type :json
  res.body
end


get '/tags' do
  url = URI("#{load_balancer.next_item}/getTags")
  if url.nil?
    return 404
  end

  req = Net::HTTP::Post.new(url, initheader = {'Content-Type' =>'application/json'})

  res = Net::HTTP.start(url.hostname, url.port) do |http|
    http.request(req)
  end

  return res.code.to_i if res.code != '200'

  content_type :json
  res.body
end