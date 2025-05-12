import type { LoaderFunctionArgs, MetaFunction } from "@remix-run/node";
import { json } from "@remix-run/node";
import { Form, useLoaderData, useSubmit } from "@remix-run/react";
import { useEffect, useRef, useCallback, useState } from "react";
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
  
  // 添加下拉框状态
  const [bvidDropdownOpen, setBvidDropdownOpen] = useState(false);
  const [intervalDropdownOpen, setIntervalDropdownOpen] = useState(false);
  const [selectedBvid, setSelectedBvid] = useState(bvid);
  const [selectedInterval, setSelectedInterval] = useState(interval);
  
  // 下拉框引用，用于点击外部关闭
  const bvidDropdownRef = useRef<HTMLDivElement>(null);
  const intervalDropdownRef = useRef<HTMLDivElement>(null);
  
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

  // 点击外部关闭下拉框
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (bvidDropdownRef.current && !bvidDropdownRef.current.contains(event.target as Node)) {
        setBvidDropdownOpen(false);
      }
      if (intervalDropdownRef.current && !intervalDropdownRef.current.contains(event.target as Node)) {
        setIntervalDropdownOpen(false);
      }
    }
    
    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  // 选择BVID
  const handleBvidSelect = (bvidOption: string) => {
    setSelectedBvid(bvidOption);
    setBvidDropdownOpen(false);
  };

  // 选择时间间隔
  const handleIntervalSelect = (intervalOption: string) => {
    setSelectedInterval(intervalOption);
    setIntervalDropdownOpen(false);
  };

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
          {/* BVID 自定义下拉框 */}
          <div>
            <label htmlFor="bvid" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">BVID:</label>
            <div className="relative" ref={bvidDropdownRef}>
              <input 
                type="hidden" 
                name="bvid" 
                value={selectedBvid}
              />
              <button
                type="button"
                className="relative w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg shadow-sm pl-3 pr-10 py-2.5 text-left cursor-default focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-sm"
                onClick={() => setBvidDropdownOpen(!bvidDropdownOpen)}
                aria-haspopup="listbox"
                aria-expanded={bvidDropdownOpen}
              >
                <span className="block truncate text-gray-900 dark:text-gray-100">{selectedBvid}</span>
                <span className="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
                  <svg className="h-5 w-5 text-gray-400 dark:text-gray-500" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                    <path fillRule="evenodd" d="M10 3a.75.75 0 01.55.24l3.25 3.5a.75.75 0 11-1.1 1.02L10 4.852 7.3 7.76a.75.75 0 01-1.1-1.02l3.25-3.5A.75.75 0 0110 3zm-3.76 9.2a.75.75 0 011.06.04l2.7 2.908 2.7-2.908a.75.75 0 111.1 1.02l-3.25 3.5a.75.75 0 01-1.1 0l-3.25-3.5a.75.75 0 01.04-1.06z" clipRule="evenodd" />
                  </svg>
                </span>
              </button>

              {bvidDropdownOpen && (
                <div className="absolute z-10 mt-1 w-full bg-white dark:bg-gray-800 shadow-lg max-h-60 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none sm:text-sm">
                  <ul tabIndex={-1} role="listbox" aria-labelledby="bvid-dropdown">
                    {bvidList && bvidList.length > 0 ? (
                      bvidList.map((bvidOption, index) => (
                        <li
                          key={index}
                          className={`cursor-default select-none relative py-2 pl-3 pr-9 hover:bg-blue-100 dark:hover:bg-blue-900 ${
                            bvidOption === selectedBvid ? 'bg-blue-50 dark:bg-blue-800 text-blue-700 dark:text-blue-300' : 'text-gray-900 dark:text-gray-100'
                          }`}
                          onClick={() => handleBvidSelect(bvidOption)}
                          role="option"
                          aria-selected={bvidOption === selectedBvid}
                        >
                          <span className="block truncate font-medium">{bvidOption}</span>
                          {bvidOption === selectedBvid && (
                            <span className="absolute inset-y-0 right-0 flex items-center pr-4 text-blue-600 dark:text-blue-400">
                              <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                                <path fillRule="evenodd" d="M16.704 4.153a.75.75 0 01.143 1.052l-8 10.5a.75.75 0 01-1.127.075l-4.5-4.5a.75.75 0 011.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 011.05-.143z" clipRule="evenodd" />
                              </svg>
                            </span>
                          )}
                        </li>
                      ))
                    ) : (
                      <li
                        className="cursor-default select-none relative py-2 pl-3 pr-9 text-gray-900 dark:text-gray-100 bg-blue-50 dark:bg-blue-800"
                        onClick={() => handleBvidSelect(bvid)}
                      >
                        <span className="block truncate font-medium">{bvid}</span>
                      </li>
                    )}
                  </ul>
                </div>
              )}
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
          
          {/* Interval 自定义下拉框 */}
          <div>
            <label htmlFor="interval" className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">Interval:</label>
            <div className="relative" ref={intervalDropdownRef}>
              <input 
                type="hidden" 
                name="interval" 
                value={selectedInterval}
              />
              <button
                type="button"
                className="relative w-full bg-white dark:bg-gray-800 border border-gray-300 dark:border-gray-600 rounded-lg shadow-sm pl-3 pr-10 py-2.5 text-left cursor-default focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 text-sm"
                onClick={() => setIntervalDropdownOpen(!intervalDropdownOpen)}
                aria-haspopup="listbox"
                aria-expanded={intervalDropdownOpen}
              >
                <span className="block truncate text-gray-900 dark:text-gray-100">
                  {selectedInterval === '10m' ? '10 Minutes' : 
                   selectedInterval === '30m' ? '30 Minutes' : 
                   selectedInterval === '1h' ? '1 Hour' : '1 Day'}
                </span>
                <span className="absolute inset-y-0 right-0 flex items-center pr-2 pointer-events-none">
                  <svg className="h-5 w-5 text-gray-400 dark:text-gray-500" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                    <path fillRule="evenodd" d="M10 3a.75.75 0 01.55.24l3.25 3.5a.75.75 0 11-1.1 1.02L10 4.852 7.3 7.76a.75.75 0 01-1.1-1.02l3.25-3.5A.75.75 0 0110 3zm-3.76 9.2a.75.75 0 011.06.04l2.7 2.908 2.7-2.908a.75.75 0 111.1 1.02l-3.25 3.5a.75.75 0 01-1.1 0l-3.25-3.5a.75.75 0 01.04-1.06z" clipRule="evenodd" />
                  </svg>
                </span>
              </button>

              {intervalDropdownOpen && (
                <div className="absolute z-10 mt-1 w-full bg-white dark:bg-gray-800 shadow-lg max-h-60 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none sm:text-sm">
                  <ul tabIndex={-1} role="listbox" aria-labelledby="interval-dropdown">
                    <li
                      className={`cursor-default select-none relative py-2 pl-3 pr-9 hover:bg-blue-100 dark:hover:bg-blue-900 ${
                        selectedInterval === '10m' ? 'bg-blue-50 dark:bg-blue-800 text-blue-700 dark:text-blue-300' : 'text-gray-900 dark:text-gray-100'
                      }`}
                      onClick={() => handleIntervalSelect('10m')}
                      role="option"
                      aria-selected={selectedInterval === '10m'}
                    >
                      <span className="block truncate font-medium">10 Minutes</span>
                      {selectedInterval === '10m' && (
                        <span className="absolute inset-y-0 right-0 flex items-center pr-4 text-blue-600 dark:text-blue-400">
                          <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                            <path fillRule="evenodd" d="M16.704 4.153a.75.75 0 01.143 1.052l-8 10.5a.75.75 0 01-1.127.075l-4.5-4.5a.75.75 0 011.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 011.05-.143z" clipRule="evenodd" />
                          </svg>
                        </span>
                      )}
                    </li>
                    <li
                      className={`cursor-default select-none relative py-2 pl-3 pr-9 hover:bg-blue-100 dark:hover:bg-blue-900 ${
                        selectedInterval === '30m' ? 'bg-blue-50 dark:bg-blue-800 text-blue-700 dark:text-blue-300' : 'text-gray-900 dark:text-gray-100'
                      }`}
                      onClick={() => handleIntervalSelect('30m')}
                      role="option"
                      aria-selected={selectedInterval === '30m'}
                    >
                      <span className="block truncate font-medium">30 Minutes</span>
                      {selectedInterval === '30m' && (
                        <span className="absolute inset-y-0 right-0 flex items-center pr-4 text-blue-600 dark:text-blue-400">
                          <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                            <path fillRule="evenodd" d="M16.704 4.153a.75.75 0 01.143 1.052l-8 10.5a.75.75 0 01-1.127.075l-4.5-4.5a.75.75 0 011.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 011.05-.143z" clipRule="evenodd" />
                          </svg>
                        </span>
                      )}
                    </li>
                    <li
                      className={`cursor-default select-none relative py-2 pl-3 pr-9 hover:bg-blue-100 dark:hover:bg-blue-900 ${
                        selectedInterval === '1h' ? 'bg-blue-50 dark:bg-blue-800 text-blue-700 dark:text-blue-300' : 'text-gray-900 dark:text-gray-100'
                      }`}
                      onClick={() => handleIntervalSelect('1h')}
                      role="option"
                      aria-selected={selectedInterval === '1h'}
                    >
                      <span className="block truncate font-medium">1 Hour</span>
                      {selectedInterval === '1h' && (
                        <span className="absolute inset-y-0 right-0 flex items-center pr-4 text-blue-600 dark:text-blue-400">
                          <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                            <path fillRule="evenodd" d="M16.704 4.153a.75.75 0 01.143 1.052l-8 10.5a.75.75 0 01-1.127.075l-4.5-4.5a.75.75 0 011.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 011.05-.143z" clipRule="evenodd" />
                          </svg>
                        </span>
                      )}
                    </li>
                    <li
                      className={`cursor-default select-none relative py-2 pl-3 pr-9 hover:bg-blue-100 dark:hover:bg-blue-900 ${
                        selectedInterval === '1d' ? 'bg-blue-50 dark:bg-blue-800 text-blue-700 dark:text-blue-300' : 'text-gray-900 dark:text-gray-100'
                      }`}
                      onClick={() => handleIntervalSelect('1d')}
                      role="option"
                      aria-selected={selectedInterval === '1d'}
                    >
                      <span className="block truncate font-medium">1 Day</span>
                      {selectedInterval === '1d' && (
                        <span className="absolute inset-y-0 right-0 flex items-center pr-4 text-blue-600 dark:text-blue-400">
                          <svg className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                            <path fillRule="evenodd" d="M16.704 4.153a.75.75 0 01.143 1.052l-8 10.5a.75.75 0 01-1.127.075l-4.5-4.5a.75.75 0 011.06-1.06l3.894 3.893 7.48-9.817a.75.75 0 011.05-.143z" clipRule="evenodd" />
                          </svg>
                        </span>
                      )}
                    </li>
                  </ul>
                </div>
              )}
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
