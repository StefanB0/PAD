require_relative '../lib/service_discovery.rb'

class LoadBalancer
  def initialize(name)
    @name = name
    @service_discovery = ServiceDiscovery.instance
    @logger = Logger.new(STDERR)

    @shift = 0
  end

  def query_service_discovery
    @items = []
    svc_address = @service_discovery.get_service_address(@name)
    unless svc_address.nil? || svc_address.empty?
      items = JSON.parse(svc_address)['services']
      self.add_item(items)
    end

    (0..@shift).each do |i|
      @items << @items.shift
    end
  end

  def add_item(item)
    @items ||= []
    if item.is_a?(Array)
      @items.concat(item)
    else
    @items << item
    end

    @items.flatten!
    @items.uniq!
  end

  def remove_item(item)
    @items.delete(item)
  end

  def next_item
    self.query_service_discovery
    @shift += 1
    @shift %= @items.length
    (@items << @items.shift).last
  end

  def log_items
    @logger.info(@items)
  end
end
