import axios from "axios";
import { MqttClient } from "mqtt";
import { json } from "stream/consumers";
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
    getData: <T> (code: string, host: string, port: number | null, queueName: string) => Promise<T | null> | null;
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
    getData: async <T> (code: string, host: string, port: number | null, queueName: string): Promise<T | null> => {
        const decode = (str: string): string => Buffer.from(str, 'base64').toString('binary');
        
        const mqtt = require('mqtt/dist/mqtt');
        const url = "http://" + host + ":" + port + "/ws";
        
        const client = mqtt.connect(url) as MqttClient;

        let data = null;

        client.subscribe(queueName, () => {
            console.log("Subscribed to " + queueName);
            client.on('message', (topic, message) => {
                data = JSON.parse(decode(message.toLocaleString())) as T
            });
        });

        return Promise.resolve(data);
    }
}
