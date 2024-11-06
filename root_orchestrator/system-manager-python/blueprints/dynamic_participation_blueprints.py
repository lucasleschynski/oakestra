import logging

from flask import request
from flask.views import MethodView
from flask_smorest import Blueprint, abort

from resource_abstractor_client import cluster_operations

dynamic_join_bp = Blueprint("Dynamic Joining", "enabling dynamic participation", url_prefix="/api/dynamic")

worker_join_schema = {
    "type": "object",
    "properties": {
        "worker_ip": {"type": "string"},
        # "worker_location": {"type": "string"},
    },
}

@dynamic_join_bp.route("/register_intent")
class DynamicJoinController(MethodView):
    @dynamic_join_bp.arguments(schema=worker_join_schema, location="json", validate=False, unknown=True)
    def post(self, *args, **kwargs):
        data = request.get_json()
        logging.log(logging.INFO, data)
        worker_ip = data.get("worker_ip")

        
        clusters = cluster_operations.get_resources(active=True)
        if clusters is None:
            return abort(500, "Getting clusters failed")
        
        logging.log(clusters)
        # return json_util.dumps(clusters)

        return "ok"
