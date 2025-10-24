import asyncio
import logging
from connectrpc import (
    request,
)
from plugin.v1 import (
    plugin_connect,
    plugin_pb2,
)
from admin.v1 import admin_connect, admin_pb2

HOST_ADDR = "http://localhost:8087"
PLUGIN_ADDR = "localhost:50102"
PLUGIN_ID = "EXTERNAL_PLUGIN_PY"

logging.basicConfig(level=logging.INFO)


class GreeterPlugin(plugin_connect.PluginService):
    async def greet(
        self,
        req: plugin_pb2.GreetRequest,
        ctx: request.RequestContext,
    ) -> plugin_pb2.GreetResponse:
        message = f"مرحبا، {req.name}"
        logging.info(
            "Responding to greet(%s): %s",
            req.name,
            message,
        )
        return plugin_pb2.GreetResponse(message=message)

    async def stream_greet(
        self,
        request: plugin_pb2.StreamGreetRequest,
        ctx,
    ):
        name = request.name
        for i in range(5):
            yield plugin_pb2.StreamGreetResponse(
                message=f"مرحبا #{i + 1} to {name} from external Python plugin"
            )


app = plugin_connect.PluginServiceASGIApplication(GreeterPlugin())


# --- Register plugin with host ---
async def register_plugin():
    client = admin_connect.AdminServiceClient(HOST_ADDR)

    logging.info("Registering plugin %s with host...", PLUGIN_ID)
    request = admin_pb2.RegisterPluginRequest(
        id=PLUGIN_ID,
        address=PLUGIN_ADDR,
    )

    try:
        response = await client.register_plugin(request)
        logging.info(
            "Registered successfully: %s",
            response.status if hasattr(response, "status") else "ok",
        )
    except Exception as e:
        logging.error(
            "Failed to register plugin %s: %s [%s]",
            PLUGIN_ID,
            e.message,
            e.code.name,
        )
        raise e


async def main():
    await register_plugin()

    logging.info(
        "Starting plugin server on %s",
        PLUGIN_ADDR,
    )

    import uvicorn

    (
        host,
        port,
    ) = PLUGIN_ADDR.split(":")
    config = uvicorn.Config(
        app,
        host=host,
        port=int(port),
        log_level="info",
    )
    server = uvicorn.Server(config)
    await server.serve()


if __name__ == "__main__":
    asyncio.run(main())
