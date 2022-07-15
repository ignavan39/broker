import axios from "axios";
import mqtt from "mqtt";
import { Host } from "../config";

export type GetInvitationResponse = {
    id: string,
    workspace: {
        id: string,
        name: string,
        createdAt: Date,
        isPrivate: boolean,
    }
    senderId: string,
    status: string, 
    systemStatus: string,
    createdAt: string,
}

export type ConnectResponse = {
    queueName: string,
    exchangeName: string,
    host: string,
    port: number,
    user: string,
    vhost: string,
    password: string,
}

export type ConnectionService = {
    connect: () => Promise<ConnectResponse>;
    getData: (code: string, host: string, port: number, queueName: string) => Promise<void> | null;
}

export const connectionService: ConnectionService = {
    connect: async (): Promise<ConnectResponse> => {
        const url = Host + "/invitations/connect";
        const apiResponse = await axios.post<ConnectResponse>(url);

        return apiResponse.data
    },
    getData: async (code: string, host: string, port: number, queueName: string): Promise<void> => {
        const url = host + ":" + port + "/ws";

        const options = {
            clean: true, 
        }

        const client = mqtt.connect(url, options)

        client.on('connect', () => {
            console.log("Connected");
            client.subscribe(queueName, () => {
                console.log("Subscribed");
            })
        });
    }
}
