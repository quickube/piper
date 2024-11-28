#enable the Admin->Network->Outbound requests -> allow for webhooks
ApplicationSetting.current.update!(allow_local_requests_from_web_hooks_and_services: true)

#create new group
Group.create!(name: "Pied Pipers", path: "pied-pipers")
g = Group.find_by(name: "Pied Pipers")
n = Namespace.find_by(name: "Pied Pipers")

#create new user
u = User.new(username: 'piper-user', email: 'piper@example.com', name: 'piper-user', password: 'Aa123456', password_confirmation: 'Aa123456')
u.assign_personal_namespace(g.organization)
u.skip_confirmation! # Use only if you want the user to be automatically confirmed. If you do not use this, the user receives a confirmation email.
u.save!

# create user token
token = u.personal_access_tokens.create(scopes: [:read_api, :write_repository, :api], name: 'p-token', expires_at: 365.days.from_now)
utoken = token.token
token.save!

#add user to group
g.add_member(u, :maintainer)
g.save!

#create new project
project = g.projects.create(name: "piper-e2e-test", path: "piper-e2e-test", creator:u)
project.save!
g.save!



#GROUP ACCESS TOKEN:
# Create the group bot user. For further group access tokens, the username should be `group_{group_id}_bot_{random_string}` and email address `group_{group_id}_bot_{random_string}@noreply.{Gitlab.config.gitlab.host}`.
admin = User.find(1)
random_string = SecureRandom.hex(16)
bot = Users::CreateService.new(admin, {name: 'g_token', username: "group_#{g.id}_bot_#{random_string}", email: "group_#{g.id}_bot_#{random_string}@noreply.gitlab.local", user_type: :project_bot }).execute

# Confirm the group bot.
bot.confirm

# Add the bot to the group with the required role.
g.add_member(bot, :maintainer)
token = bot.personal_access_tokens.create(scopes:[:read_api, :write_repository, :api], name: 'g-token', expires_at: 365.days.from_now)

# Get the token value.
gtoken = token.token

puts "USER_TOKEN #{utoken}"
puts "GROUP_TOKEN #{gtoken}"
