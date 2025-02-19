"use strict";
var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
var __rest = (this && this.__rest) || function (s, e) {
    var t = {};
    for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p) && e.indexOf(p) < 0)
        t[p] = s[p];
    if (s != null && typeof Object.getOwnPropertySymbols === "function")
        for (var i = 0, p = Object.getOwnPropertySymbols(s); i < p.length; i++) {
            if (e.indexOf(p[i]) < 0 && Object.prototype.propertyIsEnumerable.call(s, p[i]))
                t[p[i]] = s[p[i]];
        }
    return t;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.getCourseAvailability = exports.updateAvailability = exports.getTutorAvailability = void 0;
const client_1 = require("@prisma/client");
const date_fns_1 = require("date-fns");
const prisma = new client_1.PrismaClient();
const getTutorAvailability = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const { userId } = req.body;
        const tutor = yield prisma.tutor.findUnique({
            where: {
                id: userId,
            },
            include: {
                availability: {
                    include: { days: true },
                },
            },
        });
        if (!tutor || !tutor.availability) {
            res.json({
                message: "Tutor not found",
                data: null,
            });
            return;
        }
        // Transform the availability data into the format expected by the form
        const availabilityData = {
            timeGap: tutor.availability.timeGap,
        };
        [
            "monday",
            "tuesday",
            "wednesday",
            "thursday",
            "friday",
            "saturday",
            "sunday",
        ].forEach((day) => {
            var _a;
            const dayAvailability = (_a = tutor.availability) === null || _a === void 0 ? void 0 : _a.days.find((d) => d.day === day.toUpperCase());
            availabilityData[day] = {
                isAvailable: !!dayAvailability,
                startTime: dayAvailability
                    ? dayAvailability.startTime.toISOString().slice(11, 16)
                    : "09:00",
                endTime: dayAvailability
                    ? dayAvailability.endTime.toISOString().slice(11, 16)
                    : "17:00",
            };
        });
        res.json({
            message: "Availability retrieved successfully",
            data: availabilityData,
        });
    }
    catch (error) {
        res.status(500).json({ message: "Error retrieving availability", error });
    }
});
exports.getTutorAvailability = getTutorAvailability;
const updateAvailability = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    try {
        const _a = req.body, { userId } = _a, data = __rest(_a, ["userId"]);
        const tutor = yield prisma.tutor.findUnique({
            where: {
                id: userId,
            },
            include: {
                availability: true,
            },
        });
        if (!tutor) {
            res.json({
                message: "Tutor not found",
                data: null,
            });
            return;
        }
        const availabilityData = Object.entries(data).flatMap(([day, { isAvailable, startTime, endTime }]) => {
            if (isAvailable) {
                const baseDate = new Date().toISOString().split("T")[0];
                return [
                    {
                        day: day.toUpperCase(),
                        startTime: new Date(`${baseDate}T${startTime}:00Z`),
                        endTime: new Date(`${baseDate}T${endTime}:00Z`),
                    },
                ];
            }
            return [];
        });
        let updatedAvailability;
        if (tutor.availability) {
            updatedAvailability = yield prisma.availability.update({
                where: {
                    id: tutor.availability.id,
                },
                data: {
                    timeGap: data.timeGap,
                    days: {
                        deleteMany: {},
                        create: availabilityData,
                    },
                },
            });
        }
        else {
            updatedAvailability = yield prisma.availability.create({
                data: {
                    tutor_id: tutor.id,
                    timeGap: data.timeGap,
                    days: {
                        create: availabilityData,
                    },
                },
            });
        }
        res.json({
            message: "Availability updated successfully",
            data: updatedAvailability,
        });
    }
    catch (error) {
        res.status(500).json({ message: "Error retrieving availability", error });
    }
});
exports.updateAvailability = updateAvailability;
const getCourseAvailability = (req, res) => __awaiter(void 0, void 0, void 0, function* () {
    var _a, _b;
    try {
        const { courseId } = req.params;
        const course = yield prisma.course.findUnique({
            where: {
                id: Number(courseId),
            },
            include: {
                tutor: {
                    select: {
                        availability: {
                            select: {
                                days: true,
                                timeGap: true,
                            },
                        },
                    },
                },
                courseSubscriptions: {
                    select: {
                        teachingSessions: {
                            select: {
                                startTime: true,
                                endTime: true,
                            },
                        },
                    },
                },
            },
        });
        if (!course) {
            res.status(404).json({ message: "Course not found" });
            return;
        }
        const startDate = (0, date_fns_1.startOfDay)(new Date().getTime() + 7 * 60 * 60 * 1000);
        const endDate = (0, date_fns_1.addDays)(startDate, 7);
        const sessions = course.courseSubscriptions.flatMap((sub) => sub.teachingSessions);
        const availableDates = [];
        for (let date = startDate; date <= endDate; date = (0, date_fns_1.addDays)(date, 1)) {
            const dayOfWeek = (0, date_fns_1.format)(date, "EEEE").toUpperCase();
            const dayAvailability = (_a = course.tutor.availability) === null || _a === void 0 ? void 0 : _a.days.find((d) => d.day === dayOfWeek);
            if (dayAvailability) {
                const dateStr = (0, date_fns_1.format)(date, "yyyy-MM-dd");
                const slots = generateAvailableTimeSlots({
                    startTime: dayAvailability.startTime,
                    endTime: dayAvailability.endTime,
                    sessions,
                    dateStr,
                    timeGap: (_b = course.tutor.availability) === null || _b === void 0 ? void 0 : _b.timeGap,
                });
                availableDates.push({
                    date: dateStr,
                    slots,
                });
            }
        }
        res.json({
            message: "Course Availability retrieved successfully",
            data: { availableDates, course },
        });
    }
    catch (error) {
        res.status(500).json({ message: "Error retrieving availability", error });
    }
});
exports.getCourseAvailability = getCourseAvailability;
function generateAvailableTimeSlots({ startTime, endTime, sessions, dateStr, timeGap = 10, duration = 50, // Default slot duration in minutes
 }) {
    const slots = [];
    let currentTime = (0, date_fns_1.parseISO)(`${dateStr}T${startTime.toISOString().slice(11, 16)}:00.000Z`);
    const slotEndTime = (0, date_fns_1.parseISO)(`${dateStr}T${endTime.toISOString().slice(11, 16)}:00.000Z`);
    // If the date is today, start from the next available slot after the current time
    const now = new Date(new Date().getTime() + 7 * 60 * 60 * 1000);
    if ((0, date_fns_1.format)(now, "yyyy-MM-dd") === dateStr) {
        currentTime = (0, date_fns_1.isBefore)(currentTime, now)
            ? (0, date_fns_1.addMinutes)(now, timeGap)
            : currentTime;
    }
    while (currentTime < slotEndTime) {
        const slotEnd = new Date(currentTime.getTime() + duration * 60000 + timeGap * 60000);
        const isSlotAvailable = !sessions.some((booking) => {
            const bookingStart = booking.startTime;
            const bookingEnd = booking.endTime;
            return ((currentTime >= bookingStart && currentTime < bookingEnd) ||
                (slotEnd > bookingStart && slotEnd <= bookingEnd) ||
                (currentTime <= bookingStart && slotEnd >= bookingEnd));
        });
        if (isSlotAvailable) {
            slots.push(currentTime.toISOString().slice(11, 16));
        }
        currentTime = slotEnd;
    }
    return slots;
}
