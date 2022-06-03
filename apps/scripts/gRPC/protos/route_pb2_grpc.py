# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import route_pb2 as route__pb2


class RouteGuideStub(object):
    """Interface exported by the server. argument => request by client, return => response by server
    """

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.GetFeature = channel.unary_unary(
                '/protos.RouteGuide/GetFeature',
                request_serializer=route__pb2.Point.SerializeToString,
                response_deserializer=route__pb2.Feature.FromString,
                )
        self.ListFeatures = channel.unary_stream(
                '/protos.RouteGuide/ListFeatures',
                request_serializer=route__pb2.Rectangle.SerializeToString,
                response_deserializer=route__pb2.Feature.FromString,
                )
        self.RecordRoute = channel.stream_unary(
                '/protos.RouteGuide/RecordRoute',
                request_serializer=route__pb2.Point.SerializeToString,
                response_deserializer=route__pb2.RouteSummary.FromString,
                )
        self.RouteChat = channel.stream_stream(
                '/protos.RouteGuide/RouteChat',
                request_serializer=route__pb2.RouteNote.SerializeToString,
                response_deserializer=route__pb2.RouteNote.FromString,
                )


class RouteGuideServicer(object):
    """Interface exported by the server. argument => request by client, return => response by server
    """

    def GetFeature(self, request, context):
        """simple
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ListFeatures(self, request, context):
        """streaming list
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def RecordRoute(self, request_iterator, context):
        """streaming list with summary at the end
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def RouteChat(self, request_iterator, context):
        """bidirectional
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_RouteGuideServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'GetFeature': grpc.unary_unary_rpc_method_handler(
                    servicer.GetFeature,
                    request_deserializer=route__pb2.Point.FromString,
                    response_serializer=route__pb2.Feature.SerializeToString,
            ),
            'ListFeatures': grpc.unary_stream_rpc_method_handler(
                    servicer.ListFeatures,
                    request_deserializer=route__pb2.Rectangle.FromString,
                    response_serializer=route__pb2.Feature.SerializeToString,
            ),
            'RecordRoute': grpc.stream_unary_rpc_method_handler(
                    servicer.RecordRoute,
                    request_deserializer=route__pb2.Point.FromString,
                    response_serializer=route__pb2.RouteSummary.SerializeToString,
            ),
            'RouteChat': grpc.stream_stream_rpc_method_handler(
                    servicer.RouteChat,
                    request_deserializer=route__pb2.RouteNote.FromString,
                    response_serializer=route__pb2.RouteNote.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'protos.RouteGuide', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class RouteGuide(object):
    """Interface exported by the server. argument => request by client, return => response by server
    """

    @staticmethod
    def GetFeature(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/protos.RouteGuide/GetFeature',
            route__pb2.Point.SerializeToString,
            route__pb2.Feature.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ListFeatures(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_stream(request, target, '/protos.RouteGuide/ListFeatures',
            route__pb2.Rectangle.SerializeToString,
            route__pb2.Feature.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def RecordRoute(request_iterator,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.stream_unary(request_iterator, target, '/protos.RouteGuide/RecordRoute',
            route__pb2.Point.SerializeToString,
            route__pb2.RouteSummary.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def RouteChat(request_iterator,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.stream_stream(request_iterator, target, '/protos.RouteGuide/RouteChat',
            route__pb2.RouteNote.SerializeToString,
            route__pb2.RouteNote.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)