# create admin token with sudo scope to impersonate users
user = User.find_by(username: 'root');
token = user.personal_access_tokens.create!(name: 'root-token', scopes: ['sudo','api'], expires_at: 365.days.from_now); 

puts "TOKEN: #{token.token}"



#enable the Admin->Network->Outbound requests -> allow for webhooks
ApplicationSetting.current.update!(allow_local_requests_from_web_hooks_and_services: true)
