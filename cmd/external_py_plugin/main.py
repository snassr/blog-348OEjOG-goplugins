import asyncio
import logging
import httpx
from connectrpc import request
from plugin.v1 import plugin_connect, plugin_pb2
from admin.v1 import admin_connect, admin_pb2

# --- Config ---
HOST_ADDR = "http://localhost:8087"
PLUGIN_ADDR = "localhost:50102"
PLUGIN_ID = "EXTERNAL_PLUGIN_PY"

logging.basicConfig(level=logging.INFO)

# --- Implement the generated Connect service ---
class GreeterPlugin(plugin_connect.PluginService):
    async def greet(
        self, req: plugin_pb2.GreetRequest, ctx: request.RequestContext, 
    ) -> plugin_pb2.GreetResponse:
        message = f"مرحبا، {req.name}"
        logging.info("Responding to greet(%s): %s", req.name, message)
        return plugin_pb2.GreetResponse(message=message)
    
    async def stream_greet(self, request: plugin_pb2.StreamGreetRequest, ctx):
        name = request.name
        for i in range(5):
            yield plugin_pb2.StreamGreetResponse(
                message=f"مرحبا #{i+1} to {name} from external Python plugin"
            )


# --- Build the ASGI app directly from Connect ---
app = plugin_connect.PluginServiceASGIApplication(GreeterPlugin())

# --- Register plugin with host ---
async def register_plugin():
    url = f"{HOST_ADDR}/admin.v1.AdminService/RegisterPlugin"
    payload = {"id": PLUGIN_ID, "address": PLUGIN_ADDR}
    async with httpx.AsyncClient(timeout=httpx.Timeout(30.0)) as client:
        logging.info("Registering plugin %s with host...", PLUGIN_ID)
        resp = await client.post(url, json=payload)
        if resp.status_code != 200:
            logging.error("Failed to register: %s", resp.text)
        else:
            logging.info("Registered successfully: %s", resp.text)


# --- Main entrypoint ---
async def main():
    await register_plugin()
    logging.info("Starting plugin server on %s", PLUGIN_ADDR)

    import uvicorn
    host, port = PLUGIN_ADDR.split(":")
    config = uvicorn.Config(app, host=host, port=int(port), log_level="info")
    server = uvicorn.Server(config)
    await server.serve()


if __name__ == "__main__":
    asyncio.run(main())
