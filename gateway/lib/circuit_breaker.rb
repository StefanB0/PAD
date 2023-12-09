require 'logger'
require 'httparty'

# require_relative '../lib/load_balancer.rb'

class CircuitBreaker
  def initialize(service)
    @logger = Logger.new(STDERR)
    
    @load_balancer = LoadBalancer.new(service)
    
    @state = :closed
    @redirect_count = 0
    @failure_threshold = 3
    @reset_timeout = 11
    @restart_timeout = 30
    reset
  end

  def reset
    @state = :closed
    @redirect_count = 0
    @last_failure = Time.now - @reset_timeout
  end

  def open
    @state = :open
    @redirect_count = 0
    @last_failure = Time.now
  end

  def timed_reset
    if @state == :closed && (Time.now - @last_failure) > @reset_timeout
      reset
    elsif @state == :open && (Time.now - @last_failure) > @restart_timeout
      reset
    end
  end

  def call(request)
    timed_reset

    unless @state == :closed
      @logger.info("Circuit breaker is #{@state}")
    end
    
    while @redirect_count < @failure_threshold
      url = URI(@load_balancer.next_item)
      
      @logger.info("Attempting to connect to #{url}")

      if url.nil?
        return 404
      end

      begin
        if @redirect_count > 0
          @logger.info("Redirecting to #{url}")
        end

        http_client = Net::HTTP.new(url.host, url.port)
        http_client.read_timeout = 1

        res = http_client.request(request)

        break
      rescue Timeout::Error => e
        @logger.info("Connection timed out")
        @redirect_count += 1
        @last_failure = Time.now
      rescue Exception => e
        @logger.info("Connection failed with error #{e}")
        @redirect_count += 1
        @last_failure = Time.now
      end
    end

    if res.nil?
      @logger.info("Circuit breaker is open, redirect count: #{@redirect_count}")
      open
      
      return 503
    end

    return res
  end

  def call_httparty(path, method, body, headers)
    timed_reset

    unless @state == :closed
      @logger.info("Circuit breaker is #{@state}")
    end
    
    
    while @redirect_count < @failure_threshold
      url = "#{@load_balancer.next_item}#{path}"
      
      if url.nil?
        return 404
      end
      
      unless @redirect_count > 0
        @logger.info("Attempting to connect to #{url}")
      else
        @logger.info("Redirecting to #{url}")
      end
       
      begin

        req = HTTParty::Request.new(method, url, headers: headers, body: body, timeout: 1)
        res = req.perform

        break
      rescue Timeout::Error => e
        @logger.info("Connection timed out")
        @redirect_count += 1
        @last_failure = Time.now
      rescue Exception => e
        @logger.info("Connection failed with error #{e}")
        @redirect_count += 1
        @last_failure = Time.now
      end
    end

    if res.nil?
      @logger.info("Circuit breaker is open, redirect count: #{@redirect_count}")
      open
      
      return 503
    end

    return res
  end
end