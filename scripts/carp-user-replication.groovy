import groovy.json.JsonSlurper
import org.sonatype.nexus.security.user.UserNotFoundException
import org.sonatype.nexus.security.role.*
import org.apache.commons.lang.*

// parse json formatted carp user, which is send as argument for the script
def carpUser = new JsonSlurper().parseText(args)

// use undocumented getSecuritySystem to check and update existing users
def securitySystem = security.getSecuritySystem()

// default role for new users
def defaultRole = ['cesUser']

try {
  log.info('update user ' + carpUser.Username)
  def user = securitySystem.getUser(carpUser.Username)
  user.setFirstName(carpUser.FirstName)
  user.setLastName(carpUser.LastName)
  user.setEmailAddress(carpUser.Email)
  // active status and password are not changed
  securitySystem.updateUser(user)
} catch (UserNotFoundException ex) {
  log.info('create user ' + carpUser.username)

  // user not found, create a new one
  // id, firstName, lastName, Email, active, password, arrayOfRoles
  String randomUserPassword = org.apache.commons.lang.RandomStringUtils.random(16, true, true)
  security.addUser(carpUser.Username, carpUser.FirstName, carpUser.LastName, carpUser.Email, true, randomUserPassword, defaultRole)
}

// map groups to nexus roles
def authorizationManager = securitySystem.getAuthorizationManager('default')
// remove user from admin group; will be added again, if still in it
user = securitySystem.getUser(carpUser.Username)
user.removeRole(new RoleIdentifier("default", "cesAdminGroup"))
securitySystem.updateUser(user)
// add roles to user
for (group in carpUser.Groups){
  Role currentRole
  try{
    currentRole = authorizationManager.getRole(group)
  } catch (NoSuchRoleException noSuchRoleException){
    log.info('creating role ' + group)
    def newRole = new Role(
      roleId: group,
      source: "",
      name: group,
      description: "",
      readOnly: false,
      privileges: [],
      roles: []
    )
    authorizationManager.addRole(newRole)
    currentRole = newRole
  }
  user = securitySystem.getUser(carpUser.Username)
  log.info('Adding role ' + currentRole.getRoleId() + ' to user ' + user.getUserId())
  presentRole = authorizationManager.getRole(currentRole.getRoleId())
  user.addRole(new RoleIdentifier(presentRole.getSource(), presentRole.getRoleId()))
  securitySystem.updateUser(user)
}
