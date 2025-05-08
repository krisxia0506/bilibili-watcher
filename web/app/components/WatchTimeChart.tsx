import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer, TooltipProps } from 'recharts';
import { NameType, ValueType } from 'recharts/types/component/DefaultTooltipContent';

// 定义传入图表组件的数据点类型
interface ChartDataPoint {
  timeLabel: string;       // X轴显示的时间点，例如 "08:00" (基于原始时间的本地时区)
  duration: number;        // Y轴显示的观看时长（秒）
  originalStartTime: string; // 原始的 segment_start_time (UTC ISO string)
  originalEndTime: string;   // 原始的 segment_end_time (UTC ISO string)
}

// 定义 WatchSegment 类型，与 _index.tsx 中的保持一致
interface WatchSegment {
  segment_start_time: string;
  segment_end_time: string;
  watched_duration_seconds: number;
}

interface WatchTimeChartProps {
  segments: WatchSegment[];
}

// 格式化ISO日期时间字符串为本地时间的 HH:mm
const formatIsoToLocalHHMM = (isoString: string): string => {
  if (!isoString) return "N/A";
  try {
    const date = new Date(isoString); // Date 对象会自动将ISO字符串转换为本地时区
    const hours = date.getHours().toString().padStart(2, '0'); // 获取本地小时
    const minutes = date.getMinutes().toString().padStart(2, '0'); // 获取本地分钟
    return `${hours}:${minutes}`;
  } catch (e) {
    console.warn("Error formatting ISO string to local HH:MM:", isoString, e);
    return "Invalid Date";
  }
};

// 将API返回的segments数据转换为图表需要的数据格式
const transformDataForChart = (segments: WatchSegment[]): ChartDataPoint[] => {
  if (!segments || segments.length === 0) {
    return [];
  }
  return segments.map(segment => {
    return {
      timeLabel: formatIsoToLocalHHMM(segment.segment_start_time), // X轴标签使用本地时间
      duration: segment.watched_duration_seconds,
      originalStartTime: segment.segment_start_time,
      originalEndTime: segment.segment_end_time,
    };
  });
};

// 自定义 Tooltip 组件
const CustomTooltip = ({ active, payload }: TooltipProps<ValueType, NameType>) => {
  if (active && payload && payload.length) {
    const data = payload[0].payload as ChartDataPoint; // 获取当前点的数据
    
    // 将UTC ISO字符串转换为本地Date对象，以便获取本地日期和时间
    const localStartTime = new Date(data.originalStartTime);
    const localEndTime = new Date(data.originalEndTime);

    const formatDateToLocaleString = (date: Date) => {
      return date.toLocaleDateString(undefined, { // undefined 使用系统默认locale
        year: 'numeric', month: '2-digit', day: '2-digit'
      });
    };

    return (
      <div className="p-3 bg-white border border-gray-300 rounded-lg shadow-lg text-sm">
        <p className="font-semibold text-gray-700 mb-1">Segment Details</p>
        <p className="text-gray-600">
          <span className="font-medium">Start:</span> {formatIsoToLocalHHMM(data.originalStartTime)}
          <span className="text-xs text-gray-400 ml-1">({formatDateToLocaleString(localStartTime)})</span>
        </p>
        <p className="text-gray-600">
          <span className="font-medium">End:</span> {formatIsoToLocalHHMM(data.originalEndTime)}
          <span className="text-xs text-gray-400 ml-1">({formatDateToLocaleString(localEndTime)})</span>
        </p>
        <p className="text-blue-600">
          <span className="font-medium">Duration:</span> {data.duration} seconds
        </p>
      </div>
    );
  }
  return null;
};

export default function WatchTimeChart({ segments }: WatchTimeChartProps) {
  const chartData = transformDataForChart(segments);

  if (!chartData || chartData.length === 0) {
    return <p className="text-center text-gray-500 py-8">No data available to display the chart. Please adjust your filters or wait for data to be collected.</p>;
  }

  return (
    <ResponsiveContainer width="100%" height={400}>
      <LineChart
        data={chartData}
        margin={{
          top: 5,
          right: 30,
          left: 20,
          bottom: 5,
        }}
      >
        <CartesianGrid strokeDasharray="3 3" stroke="#e0e0e0" />
        <XAxis 
          dataKey="timeLabel" 
          stroke="#666"
          tick={{ fontSize: 12 }}
        />
        <YAxis 
          stroke="#666"
          tick={{ fontSize: 12 }}
          label={{ value: 'Watch Duration (seconds)', angle: -90, position: 'insideLeft', fill: '#333', dy: 40, dx: -10, fontSize: 14 }} 
        />
        <Tooltip content={<CustomTooltip />} />
        <Legend wrapperStyle={{ paddingTop: '20px' }} />
        <Line 
          type="monotone" 
          dataKey="duration" 
          name="Watched Duration"
          stroke="#3b82f6" 
          strokeWidth={2} 
          activeDot={{ r: 8, strokeWidth: 2, fill: '#3b82f6' }} 
          dot={{ r: 4, strokeWidth: 1, fill: '#fff', stroke: '#3b82f6' }}
        />
      </LineChart>
    </ResponsiveContainer>
  );
} 