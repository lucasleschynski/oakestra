import logging

from bson import json_util
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
    # @dynamic_join_bp.arguments(schema=worker_join_schema, location="json", validate=False, unknown=True)
    def post(self, *args, **kwargs):
        data = request.get_json()
        logging.log(logging.INFO, data)
        # worker_ip = data.get("worker_ip")

        clusters = cluster_operations.get_resources(active=True)
        if clusters is None:
            return abort(500, "Getting clusters failed")
        
        logging.log(logging.INFO, clusters)

        chosen_cluster = choose_cluster(clusters)

        response = {
            "cluster_manager_addr": chosen_cluster["ip"],
            "cluster_manager_port": chosen_cluster["port"]
        }
        
        return json_util.dumps(response)
    

## TODO: implement something for this
def choose_cluster(clusters):
    return clusters[0]
