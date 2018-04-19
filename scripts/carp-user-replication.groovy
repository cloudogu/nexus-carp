import groovy.json.JsonSlurper
import org.sonatype.nexus.security.user.UserNotFoundException
import org.sonatype.nexus.security.role.*

// parse json formatted carp user, which is send as argument for the script
def carpUser = new JsonSlurper().parseText(args)

// use undocumented getSecuritySystem to check and update existing users
def securitySystem = security.getSecuritySystem()

// every one should be an admin ;)
def adminRole = ['nx-admin']

try {
  log.info('update user ' + carpUser.Username)
  def user = securitySystem.getUser(carpUser.Username)
  user.setFirstName(carpUser.FirstName)
  user.setLastName(carpUser.LastName)
  user.setEmailAddress(carpUser.Email)
  // set active? password?
  securitySystem.updateUser(user)
} catch (UserNotFoundException ex) {
  log.info('create user ' + carpUser.username)

  // user not found, create a new one
  // id, firstName, lastName, Email, active, password, arrayOfRoles
  // what about the password, null is not accepted ? generate random ?
  security.addUser(carpUser.Username, carpUser.FirstName, carpUser.LastName, carpUser.Email, true, "secretPwd", adminRole)
}

// map groups to nexus roles
def authorizationManager = securitySystem.getAuthorizationManager('default')
for (group in carpUser.Groups){
  Role currentRole
  try{
    currentRole = authorizationManager.getRole(group)
  } catch (NoSuchRoleException noSuchRoleException){
    log.info('creating role ' + group)
    //TODO: add privileges if group is cesManager?
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
