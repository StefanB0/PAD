require 'sinatra'
require 'logger'

cache_manager = CacheManager.instance

logger = Logger.new(STDERR)

post '/cache' do
  node_name = params['node_name']
  node_address = params['node_address']

  logger.info "Adding node #{node_name} with address #{node_address}"
  cache_manager.add_node(node_name, node_address)
  
  201
end

delete '/cache' do
  node_name = params['node_name']
  cache_manager.remove_node(node_name)
  200
end