require 'httparty'
require 'securerandom'
require 'singleton'
require 'logger'
require 'dotenv'

require_relative '../lib/service_discovery.rb'

Dotenv.load('.env')
logger = Logger.new(STDERR)

class SagaManager
  include Singleton

  def initialize
    @saga_ledger = {}
    @service_discovery = ServiceDiscovery.instance
    @logger = Logger.new(STDERR)
  end

  def new_transaction(services)
    transaction_id = SecureRandom.uuid
    
    @saga_ledger[transaction_id] = {}
    
    services.each do |service|
      @saga_ledger[transaction_id][service] = {status: 'pending'}
    end

    @logger.info("New transaction created: #{transaction_id}")

    transaction_id
  end

  def get_transaction(transaction_id)
    @saga_ledger[transaction_id]
  end

  def check_transaction(transaction_id)
    transaction_success = true

    @saga_ledger[transaction_id].each do |service, status|
      if status[:status] == 'failure'
        transaction_success = false
        break
      end
    end

    transaction_success
  end

  def update_transaction(transaction_id, service, status)
    if @saga_ledger[transaction_id].nil?
      @logger.info("Transaction #{transaction_id} not found")
      return
    end

    @saga_ledger[transaction_id][service][:status] = status

    if status == 'failure'
      @logger.info("Transaction #{transaction_id} for service #{service} failed")
      self.revert_transaction(transaction_id)
    elsif status == 'success'
      @logger.info("Transaction #{transaction_id} for service #{service} completed successfully")
    end
  end

  def revert_transaction(transaction_id)
    @saga_ledger[transaction_id].each do |service, status|
      if status[:status] == 'success'
        addresses = @service_discovery.get_service_address(service)
        addresses = JSON.parse(addresses)

        address = addresses["services"].first
        @logger.info("Reverting transaction #{transaction_id} for service #{service} at #{address}")
        HTTParty.delete("#{address}/transaction/#{transaction_id}")
      end
    end
  end
end