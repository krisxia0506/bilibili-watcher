import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { json } from "@remix-run/node";
import { Form, useLoaderData } from "@remix-run/react";
import { formatISO, startOfDay, endOfDay } from 'date-fns';
import WatchTimeChart from "~/components/WatchTimeChart";

// 定义 API 响应体中 data.segments 数组元素的类型
interface WatchSegment {
  segment_start_time: string;
  segment_end_time: string;
  watched_duration_seconds: number;
}

// 定义 API 响应体类型
interface WatchSegmentsResponse {
  code: number;
  msg: string;
  data?: {
    segments: WatchSegment[];
  };
}

// 定义 LoaderData 类型
interface LoaderData {
  bvid: string;
  startTime: string;
  endTime: string;
  interval: string;
  segments: WatchSegment[];
  error?: string;
  // 添加一个字段来存储请求参数，方便调试时在客户端查看
  debugRequestParams?: any;
}

export const meta: MetaFunction = () => {
  return [
    { title: "Bilibili Watcher" },
    { name: "description", content: "Welcome to Bilibili Watcher!" },
  ];
};

// 后端数据获取逻辑
export async function loader({ request }: LoaderFunctionArgs) {
  console.log("[Loader] Received request:", request.url);
  const url = new URL(request.url);
  const today = new Date();

  const bvidFromParams = url.searchParams.get("bvid");
  const startTimeFromParams = url.searchParams.get("startTime");
  const endTimeFromParams = url.searchParams.get("endTime");
  const intervalFromParams = url.searchParams.get("interval");

  const bvid = bvidFromParams || "BV1rT9EYbEJa";
  const interval = intervalFromParams || "1h";

  let finalStartTime: string;
  let finalEndTime: string;

  if (startTimeFromParams) {
    // datetime-local input value is typically YYYY-MM-DDTHH:MM (local time)
    // Append seconds, parse as local Date, then convert to UTC ISO string
    finalStartTime = new Date(startTimeFromParams + ":00").toISOString();
  } else {
    // Default to start of today in UTC
    finalStartTime = startOfDay(today).toISOString();
  }

  if (endTimeFromParams) {
    // datetime-local input value is typically YYYY-MM-DDTHH:MM (local time)
    finalEndTime = new Date(endTimeFromParams + ":00").toISOString();
  } else {
    // Default to end of today in UTC
    finalEndTime = endOfDay(today).toISOString();
  }

  const apiRequestBody = {
    bvid,
    start_time: finalStartTime,
    end_time: finalEndTime,
    interval,
  };
  
  console.log("[Loader] Sending API request with body:", JSON.stringify(apiRequestBody, null, 2));

  try {
    const response = await fetch("http://localhost:8081/api/v1/video/watch-segments", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(apiRequestBody),
    });

    if (!response.ok) {
      const errorText = await response.text();
      console.error("[Loader] API request failed:", response.status, errorText);
      return json<LoaderData>({
        bvid,
        startTime: finalStartTime,
        endTime: finalEndTime,
        interval,
        segments: [],
        error: `API request failed: ${response.status} - ${errorText}`,
        debugRequestParams: apiRequestBody,
      }, { status: response.status });
    }

    const result: WatchSegmentsResponse = await response.json();
    console.log("[Loader] Received API response:", JSON.stringify(result, null, 2));

    if (result.code !== 0) {
      console.error("[Loader] API business error:", result.msg);
      return json<LoaderData>({
        bvid,
        startTime: finalStartTime,
        endTime: finalEndTime,
        interval,
        segments: [],
        error: `API error: ${result.msg}`,
        debugRequestParams: apiRequestBody,
      });
    }
    
    return json<LoaderData>({
      bvid,
      startTime: finalStartTime,
      endTime: finalEndTime,
      interval,
      segments: result.data?.segments || [],
      debugRequestParams: apiRequestBody,
    });
  } catch (error) {
    let errorMessage = "Failed to fetch watch segments. Please check the console for more details.";
    if (error instanceof Error) {
        errorMessage = error.message;
    }
    console.error("[Loader] Catch block error:", error);
    return json<LoaderData>({
      bvid,
      startTime: finalStartTime,
      endTime: finalEndTime,
      interval,
      segments: [],
      error: errorMessage,
      debugRequestParams: apiRequestBody,
    }, { status: 500 });
  }
}

// 页面组件
export default function Index() {
  const { bvid, startTime, endTime, interval, segments, error, debugRequestParams } = useLoaderData<typeof loader>();
  
  console.log("[Index Component] Received loader data:", 
    { bvid, startTime, endTime, interval, segments_count: segments?.length, error, debugRequestParams }
  );

  const formatDateTimeLocal = (isoString: string): string => {
    if (!isoString) return "";
    try {
      const date = new Date(isoString);
      if (isNaN(date.getTime())) {
        // console.warn("Invalid date string for formatDateTimeLocal, returning original:", isoString);
        return isoString;
      }
      const year = date.getFullYear();
      const month = (date.getMonth() + 1).toString().padStart(2, '0');
      const day = date.getDate().toString().padStart(2, '0');
      const hours = date.getHours().toString().padStart(2, '0');
      const minutes = date.getMinutes().toString().padStart(2, '0');
      return `${year}-${month}-${day}T${hours}:${minutes}`;
    } catch (e) {
      // console.warn("Error formatting date string for formatDateTimeLocal:", isoString, e);
      return isoString;
    }
  };

  return (
    <div className="font-sans p-4 md:p-8 bg-gray-100 min-h-screen">
      <header className="mb-10">
        <h1 className="text-3xl md:text-4xl font-bold text-center text-blue-700 tracking-tight">Bilibili Watch Time Statistics</h1>
      </header>

      {/* Optional: Display debugRequestParams for quick verification on UI */}
      {/* {debugRequestParams && (
        <pre className="bg-gray-200 p-2 rounded text-xs mb-4 overflow-auto">
          Debug Info (Sent to API): {JSON.stringify(debugRequestParams, null, 2)}
        </pre>
      )} */}

      <Form method="get" className="mb-10 p-6 bg-white rounded-xl shadow-xl space-y-6 max-w-4xl mx-auto">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-x-6 gap-y-4 items-end">
          <div>
            <label htmlFor="bvid" className="block text-sm font-medium text-gray-700 mb-1">BVID:</label>
            <input
              type="text"
              name="bvid"
              id="bvid"
              defaultValue={bvid}
              className="w-full p-2.5 border border-gray-300 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition duration-150 ease-in-out"
            />
          </div>
          <div>
            <label htmlFor="startTime" className="block text-sm font-medium text-gray-700 mb-1">Start Time:</label>
            <input
              type="datetime-local"
              name="startTime"
              id="startTime"
              defaultValue={formatDateTimeLocal(startTime)}
              className="w-full p-2.5 border border-gray-300 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition duration-150 ease-in-out"
            />
          </div>
          <div>
            <label htmlFor="endTime" className="block text-sm font-medium text-gray-700 mb-1">End Time:</label>
            <input
              type="datetime-local"
              name="endTime"
              id="endTime"
              defaultValue={formatDateTimeLocal(endTime)}
              className="w-full p-2.5 border border-gray-300 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition duration-150 ease-in-out"
            />
          </div>
          <div>
            <label htmlFor="interval" className="block text-sm font-medium text-gray-700 mb-1">Interval:</label>
            <select
              name="interval"
              id="interval"
              defaultValue={interval}
              className="w-full p-2.5 border border-gray-300 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition duration-150 ease-in-out bg-white"
            >
              <option value="10m">10 Minutes</option>
              <option value="30m">30 Minutes</option>
              <option value="1h">1 Hour</option>
              <option value="1d">1 Day</option>
            </select>
          </div>
        </div>
        <div className="text-center pt-2">
          <button
            type="submit"
            className="px-8 py-2.5 bg-blue-600 text-white font-semibold rounded-lg shadow-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-75 transition duration-150 ease-in-out"
          >
            Load Statistics
          </button>
        </div>
      </Form>

      {error && (
        <div className="mb-10 p-4 bg-red-50 text-red-700 border border-red-300 rounded-lg shadow max-w-4xl mx-auto">
          <p className="font-semibold text-lg">Error:</p>
          <p>{error}</p>
        </div>
      )}

      <div className="p-4 sm:p-6 bg-white rounded-xl shadow-xl max-w-6xl mx-auto">
        <h2 className="text-2xl font-semibold text-gray-800 mb-6 text-center">Watch Durations</h2>
        <WatchTimeChart segments={segments || []} />
      </div>

      <footer className="mt-16 text-center text-sm text-gray-600">
        <p>&copy; {new Date().getFullYear()} Bilibili Watcher. All rights reserved.</p>
      </footer>
    </div>
  );
}
