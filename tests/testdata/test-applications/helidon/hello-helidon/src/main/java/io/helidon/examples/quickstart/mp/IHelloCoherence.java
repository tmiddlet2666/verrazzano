// Copyright (c) 2022, Oracle and/or its affiliates.
// Licensed under the Universal Permissive License v 1.0 as shown at https://oss.oracle.com/licenses/upl.

package io.helidon.examples.quickstart.mp;

import javax.enterprise.context.ApplicationScoped;

import org.eclipse.microprofile.rest.client.inject.RegisterRestClient;

import javax.ws.rs.GET;
import javax.ws.rs.Produces;
import javax.ws.rs.core.MediaType;

@ApplicationScoped
//@RegisterRestClient(configKey = "HelloCoherence")
@RegisterRestClient(baseUri = "http://hello.svc.cluster.local")
public interface IHelloCoherence {

    @GET
    @Produces(MediaType.APPLICATION_JSON)
    public String greet();
}