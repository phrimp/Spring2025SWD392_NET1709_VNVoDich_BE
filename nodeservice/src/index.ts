import express from "express";
import dotenv from "dotenv";
import bodyParser from "body-parser";
import cors from "cors";
import helmet from "helmet";
import morgan from "morgan";
// ROUTE IMPORTS
import courseRoutes from "./routes/courseRoutes";
import tutorRoutes from "./routes/tutorRoutes";
import childRoutes from "./routes/childrenRoutes";
import bookingRoutes from "./routes/bookingRoutes";
import availabilityRoutes from "./routes/availabilityRoutes";
import teachingSessionRoutes from "./routes/teachingSessionRoutes";

// CONFIGURATIONS
dotenv.config();
const app = express();
app.use(express.json());
app.use(helmet());
app.use(helmet.crossOriginResourcePolicy({ policy: "cross-origin" }));
app.use(morgan("common"));
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({ extended: false }));
app.use(cors());

// ROUTES
app.get("/", (req, res) => {
  res.send("This is home route");
});

app.use("/courses", courseRoutes);
app.use("/availabilities", availabilityRoutes);
app.use("/tutors", tutorRoutes);
app.use("/childrens", childRoutes);
app.use("/bookings", bookingRoutes);
app.use("/teaching-sessions", teachingSessionRoutes);

// SERVER
const port = Number(process.env.PORT) || 3000;
app.listen(port, "0.0.0.0", () => {
  console.log(`Server running on part ${port}`);
});
