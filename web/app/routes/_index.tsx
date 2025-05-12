import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { json } from "@remix-run/node";
import { Form, useLoaderData, useSubmit } from "@remix-run/react";
import { useEffect, useRef, useCallback } from "react";
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
  bvidList: string[];
  startTime: string | null;
  endTime: string | null;
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

  // Define default BVID, preferring environment variable
  const defaultBvidStr = (typeof process !== 'undefined' && process.env.BILIBILI_BVID) 
                      ? process.env.BILIBILI_BVID 
                      : "BV1rT9EYbEJa";
                      
  // 解析多个 BVID
  const bvidList = defaultBvidStr.split(',').map(bvid => bvid.trim());
  const defaultBvid = bvidList[0] || "BV1rT9EYbEJa";

  const bvidFromParams = url.searchParams.get("bvid");
  const startTimeFromParams = url.searchParams.get("startTime"); // Expected to be UTC ISO string from client
  const endTimeFromParams = url.searchParams.get("endTime");     // Expected to be UTC ISO string from client
  const intervalFromParams = url.searchParams.get("interval");

  // Determine the final BVID to use: from params or the default
  const bvid = bvidFromParams || defaultBvid;
  const interval = intervalFromParams || "1h";

  // Now directly use params if they exist, assuming client sent UTC ISO
  const finalStartTimeForApi: string | null = startTimeFromParams || null;
  const finalEndTimeForApi: string | null = endTimeFromParams || null;
  
  // Pass back the same UTC ISO strings (or null) to client
  const startTimeForClientDisplay: string | null = finalStartTimeForApi;
  const endTimeForClientDisplay: string | null = finalEndTimeForApi;

  const apiRequestBody = (finalStartTimeForApi && finalEndTimeForApi) ? {
    bvid,
    start_time: finalStartTimeForApi,
    end_time: finalEndTimeForApi,
    interval,
  } : null;
  
  let segments: WatchSegment[] = [];
  let apiError: string | undefined = undefined;
  const backendApiBaseUrl = typeof process !== 'undefined' && process.env.BACKEND_API_URL ? process.env.BACKEND_API_URL : "http://localhost:8081";
  const apiUrl = `${backendApiBaseUrl}/api/v1/video/watch-segments`;

  if (apiRequestBody) {
    console.log(`[Loader] API Request Body for ${apiUrl}:`, JSON.stringify(apiRequestBody, null, 2));
    try {
      const response = await fetch(apiUrl, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(apiRequestBody),
      });
      if (!response.ok) {
        const errorText = await response.text();
        console.error("[Loader] API request failed:", response.status, errorText);
        apiError = `API request failed: ${response.status} - ${errorText}`;
      } else {
        const result: WatchSegmentsResponse = await response.json();
        console.log("[Loader] Received API response:", JSON.stringify(result, null, 2));
        if (result.code !== 0) {
          console.error("[Loader] API business error:", result.msg);
          apiError = `API error: ${result.msg}`;
        } else {
          segments = result.data?.segments || [];
        }
      }
    } catch (error) {
      let errorMessage = "Failed to fetch watch segments.";
      if (error instanceof Error) { errorMessage = error.message; }
      console.error("[Loader] Catch block error:", error);
      apiError = errorMessage;
    }
  } else {
    console.log("[Loader] startTime or endTime missing or invalid in URL query, skipping API call.");
  }
  
  return json<LoaderData>({
    bvid,
    bvidList,
    startTime: startTimeForClientDisplay,
    endTime: endTimeForClientDisplay,
    interval,
    segments,
    error: apiError,
    debugRequestParams: apiRequestBody, 
  });
}

// 页面组件
export default function Index() {
  const { bvid, bvidList, startTime: startTimeFromLoader, endTime: endTimeFromLoader, interval, segments, error } = useLoaderData<typeof loader>();
  const submit = useSubmit();
  const formRef = useRef<HTMLFormElement>(null);
  
  console.log("[Index Component] Received loader data:", 
    { bvid, startTime: startTimeFromLoader, endTime: endTimeFromLoader, interval, segments_count: segments?.length, error }
  );

  // Converts a Date object to 'YYYY-MM-DDTHH:MM' string for datetime-local input
  const formatLocalDateToInputString = (date: Date): string => {
    const year = date.getFullYear();
    const month = (date.getMonth() + 1).toString().padStart(2, '0');
    const day = date.getDate().toString().padStart(2, '0');
    const hours = date.getHours().toString().padStart(2, '0');
    const minutes = date.getMinutes().toString().padStart(2, '0');
    return `${year}-${month}-${day}T${hours}:${minutes}`;
  };

  // Converts a UTC ISO string (from loader, if present) to 'YYYY-MM-DDTHH:MM' local string for datetime-local input
  const formatUtcIsoToLocalInputString = (isoString: string): string => {
    if (!isoString) return "";
    try {
      const date = new Date(isoString); // Converts UTC ISO to local Date object
      if (isNaN(date.getTime())) return isoString;
      return formatLocalDateToInputString(date); // Then format this local Date object
    } catch (e) {
      return isoString;
    }
  };

  // Calculate initial values for the uncontrolled inputs
  let initialStartTimeForInput: string;
  let initialEndTimeForInput: string;
  if (startTimeFromLoader && endTimeFromLoader) {
    initialStartTimeForInput = formatUtcIsoToLocalInputString(startTimeFromLoader);
    initialEndTimeForInput = formatUtcIsoToLocalInputString(endTimeFromLoader);
  } else {
    const todayClient = new Date();
    const clientStartOfDay = new Date(todayClient.getFullYear(), todayClient.getMonth(), todayClient.getDate(), 0, 0, 0, 0);
    const clientEndOfDay = new Date(todayClient.getFullYear(), todayClient.getMonth(), todayClient.getDate(), 23, 59, 59, 999);
    initialStartTimeForInput = formatLocalDateToInputString(clientStartOfDay);
    initialEndTimeForInput = formatLocalDateToInputString(clientEndOfDay);
  }

  // Centralized submit logic with UTC conversion
  const handleSubmit = useCallback((formData: FormData) => {
    const startTimeLocal = formData.get("startTime") as string;
    const endTimeLocal = formData.get("endTime") as string;
    const bvidValue = formData.get("bvid") as string || bvid;
    const intervalValue = formData.get("interval") as string || interval;
    
    let startTimeUtcIso: string | null = null;
    let endTimeUtcIso: string | null = null;

    if (startTimeLocal) {
      try {
        // Append seconds and convert the local time string to UTC ISO
        startTimeUtcIso = new Date(startTimeLocal + ":00").toISOString();
      } catch (e) {
        console.error("Error parsing start time:", startTimeLocal, e);
        // Optionally handle error, e.g., show validation message
      }
    }
    if (endTimeLocal) {
      try {
        // Append seconds and convert the local time string to UTC ISO
        endTimeUtcIso = new Date(endTimeLocal + ":00").toISOString();
      } catch (e) {
        console.error("Error parsing end time:", endTimeLocal, e);
        // Optionally handle error
      }
    }
    
    if (startTimeUtcIso && endTimeUtcIso) {
      const params = new URLSearchParams();
      params.append("bvid", bvidValue);
      params.append("interval", intervalValue);
      params.append("startTime", startTimeUtcIso);
      params.append("endTime", endTimeUtcIso);
      
      console.log("[Index Component] Submitting with client-converted UTC params:", params.toString());
      submit(params, { method: "get", replace: false }); // Use replace: false for manual submits
    } else {
      console.warn("Invalid date/time input, submission aborted.");
      // Optionally provide user feedback
    }
  }, [submit, bvid, interval]);

  // Handle manual form submission
  const handleFormSubmit = useCallback((event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault(); // Prevent default GET submission
    const formData = new FormData(event.currentTarget);
    handleSubmit(formData);
  }, [handleSubmit]);

  // Handle automatic submission on initial load
  useEffect(() => {
    if (
      startTimeFromLoader === null &&
      endTimeFromLoader === null &&
      segments && segments.length === 0 &&
      !error &&
      formRef.current // Ensure form is mounted
    ) {
      console.log("[Index Component] Initial load: Triggering submit with client defaults.");
      // Create FormData from the current form state which has the client default values
      const initialFormData = new FormData(formRef.current);
      handleSubmit(initialFormData); // Use the centralized handler for UTC conversion
    }
    // Dependencies updated slightly
  }, [startTimeFromLoader, endTimeFromLoader, segments, error, handleSubmit]);

  return (
    <div className="font-sans p-4 md:p-8 bg-gray-100 dark:bg-gray-800 min-h-screen">
      <header className="mb-10">
        <h1 className="text-3xl md:text-4xl font-bold text-center text-blue-700 dark:text-blue-500 tracking-tight">Bilibili Watch Time Statistics</h1>
      </header>

      {/* Optional: Display debugRequestParams for quick verification on UI */}
      {/* {debugRequestParams && (
        <pre className="bg-gray-200 p-2 rounded text-xs mb-4 overflow-auto">
          Debug Info (Sent to API): {JSON.stringify(debugRequestParams, null, 2)}
        </pre>
      )} */}

      <Form ref={formRef} method="get" onSubmit={handleFormSubmit} className="mb-10 p-6 bg-white dark:bg-gray-700 rounded-xl shadow-xl space-y-6 max-w-4xl mx-auto">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-x-6 gap-y-4 items-end">
          <div>
            <label htmlFor="bvid" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">BVID:</label>
            <div className="relative">
              <select
                name="bvid"
                id="bvid"
                defaultValue={bvid}
                className="w-full p-2.5 border border-gray-300 dark:border-gray-600 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white dark:bg-gray-800 dark:text-gray-100 pr-10 appearance-none"
              >
                {bvidList && bvidList.length > 0 ? (
                  bvidList.map((bvidOption, index) => (
                    <option 
                      key={index} 
                      value={bvidOption}
                      className="bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100"
                    >
                      {bvidOption}
                    </option>
                  ))
                ) : (
                  <option value={bvid} className="bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100">{bvid}</option>
                )}
              </select>
              <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-gray-700 dark:text-gray-300">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7"></path>
                </svg>
              </div>
            </div>
          </div>
          <div>
            <label htmlFor="startTime" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Start Time:</label>
            <input
              type="datetime-local"
              name="startTime"
              id="startTime"
              defaultValue={initialStartTimeForInput}
              required
              className="w-full p-2.5 border border-gray-300 dark:border-gray-600 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition duration-150 ease-in-out bg-white dark:bg-gray-800 dark:text-gray-100"
            />
          </div>
          <div>
            <label htmlFor="endTime" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">End Time:</label>
            <input
              type="datetime-local"
              name="endTime"
              id="endTime"
              defaultValue={initialEndTimeForInput}
              required
              className="w-full p-2.5 border border-gray-300 dark:border-gray-600 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition duration-150 ease-in-out bg-white dark:bg-gray-800 dark:text-gray-100"
            />
          </div>
          <div>
            <label htmlFor="interval" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Interval:</label>
            <div className="relative">
              <select
                name="interval"
                id="interval"
                defaultValue={interval}
                className="w-full p-2.5 border border-gray-300 dark:border-gray-600 rounded-lg shadow-sm focus:ring-2 focus:ring-blue-500 focus:border-blue-500 bg-white dark:bg-gray-800 dark:text-gray-100 pr-10 appearance-none"
              >
                <option value="10m" className="bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100">10 Minutes</option>
                <option value="30m" className="bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100">30 Minutes</option>
                <option value="1h" className="bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100">1 Hour</option>
                <option value="1d" className="bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100">1 Day</option>
              </select>
              <div className="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-gray-700 dark:text-gray-300">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M19 9l-7 7-7-7"></path>
                </svg>
              </div>
            </div>
          </div>
        </div>
        <div className="text-center pt-2">
          <button
            type="submit"
            className="px-8 py-2.5 bg-blue-600 hover:bg-blue-700 dark:bg-blue-500 dark:hover:bg-blue-600 text-white font-semibold rounded-lg shadow-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-75 transition duration-150 ease-in-out"
          >
            Load Statistics
          </button>
        </div>
      </Form>

      {error && (
        <div className="mb-10 p-4 bg-red-50 dark:bg-red-900 dark:bg-opacity-30 text-red-700 dark:text-red-300 border border-red-300 dark:border-red-700 rounded-lg shadow max-w-4xl mx-auto">
          <p className="font-semibold text-lg">Error:</p>
          <p>{error}</p>
        </div>
      )}

      <div className="p-4 sm:p-6 bg-white dark:bg-gray-700 rounded-xl shadow-xl max-w-6xl mx-auto">
        <h2 className="text-2xl font-semibold text-gray-800 dark:text-gray-100 mb-6 text-center">Watch Durations</h2>
        <WatchTimeChart segments={segments || []} />
      </div>

      <footer className="mt-16 text-center text-sm text-gray-600 dark:text-gray-400">
        <p>&copy; {new Date().getFullYear()} Bilibili Watcher. All rights reserved.</p>
      </footer>
    </div>
  );
}
