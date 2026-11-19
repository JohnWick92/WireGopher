Here is the revised README in English, clarifying that WireGopher acts as an authentication proxy in front of the VPN, with the firewall intercepting unauthenticated VPN access and directing the user through the authentication flow. The exact implementation depends on the firewall platform and its capabilities.

# WireGopher

**WireGopher** is an authentication middleware for WireGuard VPN environments, designed to integrate VPN access control with external identity providers such as **Keycloak**.

It acts as an authentication proxy between the firewall, VPN clients, and the identity provider, allowing users to authenticate before being granted access to protected network resources.

## How It Works

WireGopher is designed to work alongside the firewall as part of an authentication flow.

The firewall is responsible for identifying VPN clients that have not yet completed authentication and intercepting or restricting their access accordingly. Instead of granting unrestricted network access, the firewall directs unauthenticated users to the WireGopher authentication flow.

WireGopher then handles the authentication process by communicating with the configured identity provider, such as Keycloak.

Once authentication succeeds, WireGopher updates the client's authorization state through the configured firewall integration. The firewall can then apply the appropriate access policies and grant access to authorized resources.

The general flow is:

1. A VPN client establishes a WireGuard connection.
2. The firewall identifies the client as unauthenticated.
3. The firewall restricts access to protected resources and directs the client to the WireGopher authentication flow.
4. WireGopher communicates with the configured identity provider.
5. The user completes authentication, including TOTP/2FA when configured.
6. WireGopher updates the client's authorization state through the firewall integration.
7. The firewall applies its rules and allows access to the authorized network resources.

**Important:** The interception and redirection mechanisms are firewall-dependent. The exact implementation will vary according to the firewall platform, its capabilities, and its configuration.

## Firewall Integration and Responsibilities

**Configuring and maintaining the firewall rules required for the integration is the responsibility of the infrastructure administrator.**

WireGopher does not automatically configure firewall policies or guarantee traffic enforcement. The administrator must ensure that the firewall correctly intercepts unauthenticated clients, communicates with the middleware as required, and enforces the resulting authorization state.

### Example: OPNsense Integration

For an OPNsense deployment, the administrator is expected to have the required aliases and firewall rules configured beforehand.

The following aliases are examples of the expected authorization model:

- `VPN_AUTHORIZED`: Contains VPN clients that have successfully completed authentication.
- `VPN_UNAUTHORIZED`: Contains VPN clients that have not authenticated successfully or do not currently have valid authorization.

The firewall rules must use these aliases to control access to protected network resources.

WireGopher is responsible for updating the client's authorization state according to the authentication result, while OPNsense remains responsible for enforcing the actual network access policies.

Other firewall platforms may require different integration mechanisms, aliases, policies, or traffic interception configurations.

## Authentication Requirements

A functional authentication provider must be configured and available before deploying WireGopher.

The authentication provider should support an appropriate authentication mechanism, preferably including multi-factor authentication (MFA), such as TOTP.

**Keycloak** is an example of a supported identity provider that can be configured to provide:

- User identity verification.
- Username and password authentication.
- TOTP-based two-factor authentication (2FA).
- Centralized identity and access management.

WireGopher relies on the configured identity provider to validate user credentials and determine the authentication outcome.

## Security and Deployment Considerations

- Firewall configuration and rule enforcement are the administrator's responsibility.
- Required firewall aliases, policies, and interception mechanisms must be configured beforehand.
- An identity provider must be available and properly configured.
- TOTP or another appropriate MFA mechanism should be configured in the identity provider.
- Unauthenticated clients must not have unrestricted access to protected network resources.
- Authorization state must be correctly synchronized between WireGopher and the firewall.
- A successful WireGuard handshake alone must never be considered proof of user authentication.

## Project Scope

WireGopher focuses on bridging the gap between WireGuard VPN connectivity and identity-based authentication.

Its primary responsibilities are:

- Acting as an authentication proxy in the VPN access flow.
- Integrating with external identity providers such as Keycloak.
- Handling authentication outcomes.
- Communicating authorization state to the configured firewall integration.

The firewall remains the enforcement point for network access policies.

**WireGopher authenticates. The firewall enforces.**
