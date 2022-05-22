import { bootstrap } from "../common/libs/bootstrap";
import { HttpServer } from "../common/libs/transport/http/http.server";
import config from "config";

class CalendarApp {
  public init = async () => {
    try {
      await bootstrap.database({
        type: "postgres",
        password: config.get<string>("calendar.database.password"),
        username: config.get<string>("calendar.database.user"),
        database: config.get<string>("calendar.database.name"),
        host: config.get<string>("calendar.database.host"),
        port: config.get<number>("calendar.database.port"),
      });
      const httpServer = new HttpServer();
      await httpServer.initializeRoutes("1", []).listen();
    } catch (e) {
      console.error(e);
      throw e;
    }
  };
}

(async () => {
  const app = new CalendarApp();
  await app.init();
})();
