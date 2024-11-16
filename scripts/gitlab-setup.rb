# create new user1
# u = User.new(username: 'piper-user', email: 'piper@example.com', name: 'piper-user', password: 'Aa123456', password_confirmation: 'Aa123456', )
# u.assign_personal_namespace(Organizations::Organization.default_organization)
# u.skip_confirmation! # Use only if you want the user to be automatically confirmed. If you do not use this, the user receives a confirmation email.
# u.save!

# # create another user token
# token = u.personal_access_tokens.create(scopes: ['read_api', 'write_repository', "api"], name: 'user_token', expires_at: 365.days.from_now)
# token.set_token('user-token-string')
# token.save!

user = User.find_by(username: 'root');
token = user.personal_access_tokens.create!(name: 'root-token', scopes: ['sudo','api'], expires_at: 365.days.from_now); 
# token.set_token("gitlab-root-token-key"); 
puts "TOKEN: #{token.token}"



#enable the Admin->Network->Outbound requests -> allow for webhooks
ApplicationSetting.current.update!(allow_local_requests_from_web_hooks_and_services: true)

# #create new group
# Group.create!(name: "Pied Pipers", path: "pied-pipers")
# g = Group.find_by(name: "Pied Pipers")

# g.add_member(u, Gitlab::Access::DEVELOPER)
# g.save!



# project = Projects::CreateService.new(u, name: "piper-e2e-test", path: "piper-e2e-test").execute
# # project.import_url = "https://github.com/quickube/piper-e2e-test.git"
# project.save!
# # project.import_state.start!