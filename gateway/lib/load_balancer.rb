require_relative '../lib/service_discovery.rb'

class LoadBalancer
  def initialize(name)
    @name = name
    @service_discovery = ServiceDiscovery.instance
    self.query_service_discovery
  end

  def query_service_discovery
    img_svc_address = @service_discovery.get_service_address(@name)
    unless img_svc_address.nil? || img_svc_address.empty?
      @items = JSON.parse(img_svc_address)['services']
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
