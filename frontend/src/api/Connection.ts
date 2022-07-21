import axios from "axios";
import { MqttClient } from "mqtt";
import { Host } from "../config";
import { Invitation, InvitationStatus, SystemStatus } from "../types/Invitation";

export type GetInvitationResponse = {
    id: string,
    workspace: {
        id: string,
        name: string,
        createdAt: Date,
        isPrivate: boolean,
    }
    senderId: string,
    status: InvitationStatus, 
    systemStatus: SystemStatus,
    createdAt: string,
}

export type ConnectResponse = {
    consume: {
        queueName: string,
        exchange: string,
        host: string,
        port: number,
        user: string,
        vhost: string,
        password: string,
    },
}

export type ConnectionService = {
    ping: () => Promise<void>;
    connect: () => Promise<ConnectResponse>;
    getData: (code: string, host: string, port: number | null, queueName: string) => Promise<any> | null;
}

export const connectionService: ConnectionService = {
    ping: async (): Promise<void> => {
        const url = Host + "/connect/ping";
        await axios.post(url);
    },
    connect: async (): Promise<ConnectResponse> => {
        const url = Host + "/invitations/connect";
        const apiResponse = await axios.post<ConnectResponse>(url);
        return apiResponse.data;
    },
    getData: async (code: string, host: string, port: number | null, queueName: string): Promise<any> => {
        const mqtt = require('mqtt/dist/mqtt');
        const url = "http://" + host + ":" + port + "/ws";
        
        const client = mqtt.connect(url) as MqttClient;

        client.subscribe(queueName, () => {
            console.log("Subscribed to " + queueName);
        });

        client.on('message', (topic, message) => {
            console.log(message);
        });
    }
}
