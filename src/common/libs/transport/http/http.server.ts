import express, { Application, NextFunction, Request, Response } from "express";
import * as bodyParser from "body-parser";
import cors from "cors";
import * as http from "http";
import { InternalRouter } from "../../types/router";

export class HttpServer {
  public app: Application;
  private port: number;
  private listener: http.Server;

  constructor(port?: number) {
    this.app = express();
    this.port = port ?? 8080;
  }

  public initializeRoutes(
    apiVersion: string,
    controllers: InternalRouter[]
  ): HttpServer {
    this.app.use(bodyParser.json({ limit: "10mb" }));
    this.app.use(bodyParser.urlencoded({ limit: "10mb", extended: true }));

    const options: cors.CorsOptions = {
      origin: "*",
      allowedHeaders: [
        "Origin",
        "Content-Length",
        "Content-Type",
        "Authorization",
      ],
      methods: ["*"],
    };
    this.app.use(cors(options));
    this.app.options(
      "/*",
      async (req: Request, res: Response, next: NextFunction) => {
        res.status(200).json();
      }
    );

    this.app.use(
      (err: Error, req: Request, res: Response, next: NextFunction) => {
        if (err instanceof SyntaxError) {
          if ("body" in err) {
            return res.status(400).json({
              code: 400,
              errors: [
                {
                  code: "bad json",
                },
              ],
            });
          }
        }

        return res.status(500).json({
          code: 500,
        });
      }
    );

    this.app.use((req: Request, res: Response, next: NextFunction) => {
      res.setHeader("Content-Type", "application/json");
      next();
    });

    this.app.all("/*", async (req, res) => {
      res.status(404).json({
        code: 404,
        message: "Method Not Found",
      });
    });

    controllers.forEach((controller) => {
      this.app.use(`/api/${apiVersion}/`, controller.getRouter());
    });

    return this;
  }

  public async listen() {
    this.listener = this.app.listen(this.port);
  }
}
