import axios from "axios";
import { Host } from "../config";
import { InvitationStatus, SystemStatus } from "../types/Invitation";
const mqtt = require("precompiled-mqtt");

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
    getData:(credential: ConnectResponse['consume']) => void;
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
    getData: async (credential: ConnectResponse['consume']) => {
        const {user,password,port,host,queueName} = credential
        const client = mqtt.connect({
            port,
            host,
            username: user,
            password,
        })

        let data = null;

        client.subscribe(queueName, () => {
            console.log("Subscribed to " + queueName);
            client.on('message', (topic: any, message: any) => {
                data = JSON.parse(message.buffer.toString());
                console.log(data)
            });
        });
    }
}
