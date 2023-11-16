require_relative '../lib/service_discovery.rb'

class LoadBalancer
  def initialize(name)
    @name = name
    @service_discovery = ServiceDiscovery.instance
  end

  def query_service_discovery
    svc_address = @service_discovery.get_service_address(@name)
    unless svc_address.nil? || svc_address.empty?
      @items = JSON.parse(svc_address)['services']
      self.add_item(@items)
    end
  end

  def add_item(item)
    @items ||= []
    if item.is_a?(Array)
      @items.concat(item)
    else
    @items << item
    end
    puts @items
    @items.flatten!
    @items.uniq!
  end

  def remove_item(item)
    @items.delete(item)
  end

  def next_item
    self.query_service_discovery
    (@items << @items.shift).last
  end
end
