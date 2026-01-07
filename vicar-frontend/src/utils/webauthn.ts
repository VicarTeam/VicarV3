export function decodePublicKeyCredentialCreationOptions(options: any): PublicKeyCredentialCreationOptions {
  return {
    attestation: options.attestation,
    authenticatorSelection: options.authenticatorSelection,
    challenge: arrayBufferDecode(options.challenge),
    excludeCredentials: options.excludeCredentials?.map(decodePublicKeyCredentialDescriptor),
    extensions: options.extensions,
    pubKeyCredParams: options.pubKeyCredParams,
    rp: options.rp,
    timeout: options.timeout,
    user: {
      displayName: options.user.displayName,
      id: arrayBufferDecode(options.user.id),
      name: options.user.name,
    },
  };
}

export function decodePublicKeyCredentialRequestOptions(options: any): PublicKeyCredentialRequestOptions {
  let allowCredentials: PublicKeyCredentialDescriptor[] | undefined = undefined;

  if (options.allowCredentials?.length !== 0) {
    allowCredentials = options.allowCredentials?.map(decodePublicKeyCredentialDescriptor);
  }

  return {
    allowCredentials: allowCredentials,
    challenge: arrayBufferDecode(options.challenge),
    extensions: options.extensions,
    rpId: options.rpId,
    timeout: options.timeout,
    userVerification: options.userVerification,
  };
}

export function encodeAttestationPublicKeyCredential(credential: any): any {
  const response = credential.response as any;

  let transports: AuthenticatorTransport[] | undefined;

  if (response?.getTransports !== undefined && typeof response.getTransports === 'function') {
    transports = response.getTransports();
  }

  return {
    id: credential.id,
    type: credential.type,
    rawId: arrayBufferEncode(credential.rawId),
    clientExtensionResults: credential.getClientExtensionResults(),
    response: {
      attestationObject: arrayBufferEncode(response.attestationObject),
      clientDataJSON: arrayBufferEncode(response.clientDataJSON),
    },
    transports: transports,
  };
}

export function encodeAssertionPublicKeyCredential(credential: any): any {
  const response = credential.response as AuthenticatorAssertionResponse;

  let userHandle: string;

  if (response.userHandle == null) {
    userHandle = "";
  } else {
    userHandle = arrayBufferEncode(response.userHandle)
  }

  return {
    id: credential.id,
    type: credential.type,
    rawId: arrayBufferEncode(credential.rawId),
    clientExtensionResults: credential.getClientExtensionResults(),
    response: {
      authenticatorData: arrayBufferEncode(response.authenticatorData),
      clientDataJSON: arrayBufferEncode(response.clientDataJSON),
      signature: arrayBufferEncode(response.signature),
      userHandle: userHandle,
    },
  };

}

function decodePublicKeyCredentialDescriptor(descriptor: any): PublicKeyCredentialDescriptor {
  return {
    id: arrayBufferDecode(descriptor.id),
    type: descriptor.type,
    transports: descriptor.transports,
  }
}

function arrayBufferEncode(value: ArrayBuffer): string {
  return uint8ArrayToBase64Url(new Uint8Array(value));
}

function arrayBufferDecode(value: string): ArrayBuffer {
  return base64UrlToUint8Array(value);
}

function base64UrlToUint8Array(base64Url: string): ArrayBuffer {
  let base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');

  const padding = base64.length % 4;
  if (padding === 2) {
    base64 += '==';
  } else if (padding === 3) {
    base64 += '=';
  }

  const rawBinary = atob(base64);
  const outputArray = new Uint8Array(rawBinary.length);
  for (let i = 0; i < rawBinary.length; i++) {
    outputArray[i] = rawBinary.charCodeAt(i);
  }

  // @ts-ignore
  return outputArray;
}

function uint8ArrayToBase64Url(uint8Array: Uint8Array): string {
  let binary = '';
  for (let i = 0; i < uint8Array.byteLength; i++) {
    // @ts-ignore
    binary += String.fromCharCode(uint8Array[i]);
  }

  return btoa(binary).replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '');
}
