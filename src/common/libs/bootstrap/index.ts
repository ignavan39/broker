import { Connection, createConnection, DataSourceOptions } from "typeorm";

export class Bootstrap {

  public database = async (options: DataSourceOptions): Promise<Connection> => {
    return await createConnection(options);
  };
}

export const bootstrap = new Bootstrap();
