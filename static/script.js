var chart_1 = "temp_humidity.csv";
var chart_2 = "sound_level.csv";



document.addEventListener("DOMContentLoaded", function() {

    // Your initial functions
    createLiveChart(`/static/${chart_1}`, "chart1");
    createLiveChart(`/static/${chart_2}`, "chart2");
    updateTime();
    loadContacts();
    setInterval(updateTime, 1000);
    loadThresholds();
  });


// Update clock
function updateTime() {
  const now = new Date();
  const timeStr = now.toLocaleTimeString();
  document.getElementById("time").textContent = timeStr;
}

function createLiveChart(fileName, chartId) {
  const points = 25; // number of points visible at once
  let xData = [];
  let yDataSets = []; // array of arrays for each Y column
  let lastTimestamp = null;
  let isDragging = false;
  let dragStartX = null;
  let dragStartRange = null;
  let userInteracting = false;

  function loadHistory() {
    fetch(fileName + "?nocache=" + Date.now())
      .then((res) => res.text())
      .then((text) => {
        let rows = text.trim().split("\n");
        if (rows.length < 2) return;

        // Parse headers
        let headers = rows[0].split(",");
        let yCount = headers.length - 1; // all columns except first are Y

        // Reset data
        xData = [];
        yDataSets = Array.from({ length: yCount }, () => []);

        // Fill data arrays
        rows.slice(1).forEach((row) => {
          let cols = row.split(",");
          let time = new Date(cols[0]);
          xData.push(time);

          cols.slice(1).forEach((val, idx) => {
            yDataSets[idx].push(parseFloat(val));
          });
        });

        lastTimestamp = xData[xData.length - 1];

        let startIndex = Math.max(0, xData.length - points);
        let visibleX = xData.slice(startIndex);
        let visibleYs = yDataSets.map((set) => set.slice(startIndex));

        // Generate traces dynamically
        let traces = visibleYs.map((ySet, idx) => ({
          x: visibleX,
          y: ySet,
          mode: "lines",
          name: headers[idx + 1], // column name
          line: { width: 1 },
        }));

    Plotly.newPlot(chartId, traces, {
  plot_bgcolor: "rgba(113, 111, 111, 0.3)",
  paper_bgcolor: "black",
  font: { color: "white" },
  xaxis: {
    autorange: false,
    range: [visibleX[0], visibleX[visibleX.length - 1]],
    showgrid: true,
    gridcolor: "rgba(200,200,200,0.3)",
    tickformat: "%H:%M:%S",
  },
  yaxis: {
    autorange: true,
    showgrid: true,
    gridcolor: "rgba(200,200,200,0.3)",

    // 🔹 Mobile-friendly tweaks
    tickangle: window.innerWidth < 600 ? -45 : 0,
    tickfont: { size: window.innerWidth < 600 ? 10 : 14 },
  },
  margin: {
    t: 20,
    b: 40,
    l: window.innerWidth < 600 ? 40 : 80, // dynamic left margin
    r: 10,
  },
  legend: {
    font: { size: window.innerWidth < 600 ? 10 : 12 },
    orientation: window.innerWidth < 600 ? "h" : "v", // horizontal on mobile
    x: 0,
    y: 1.1,
  },
  autosize: true,
  responsive: true,
});

      });
  }

  function fetchData() {
    fetch(fileName + "?nocache=" + Date.now())
      .then((res) => res.text())
      .then((text) => {
        let rows = text.trim().split("\n");
        if (rows.length < 2) return;

        let lastRow = rows[rows.length - 1].split(",");
        let time = new Date(lastRow[0]);
        if (lastTimestamp && time <= lastTimestamp) return;

        lastTimestamp = time;
        xData.push(time);

        lastRow.slice(1).forEach((val, idx) => {
          if (!yDataSets[idx]) yDataSets[idx] = [];
          yDataSets[idx].push(parseFloat(val));
        });

        let startIndex = Math.max(0, xData.length - points);
        let visibleX = xData.slice(startIndex);
        let visibleYs = yDataSets.map((set) => set.slice(startIndex));

        Plotly.update(chartId, {
          x: visibleYs.map(() => visibleX),
          y: visibleYs,
        });

        if (!userInteracting) {
          Plotly.relayout(chartId, {
            "xaxis.range": [visibleX[0], visibleX[visibleX.length - 1]],
          });
        }
      });
  }

  const chartEl = document.getElementById(chartId);

  // Zoom
  chartEl.addEventListener(
    "wheel",
    function (e) {
      e.preventDefault();
      userInteracting = true;

      let gd = document.getElementById(chartId);
      let layout = gd.layout;
      let currentRange = layout.xaxis.range;

      let start = new Date(currentRange[0]);
      let end = new Date(currentRange[1]);
      let rangeMs = end - start;
      let zoomFactor = 0.1 * rangeMs;

      if (e.deltaY < 0) {
        start = new Date(start.getTime() + zoomFactor);
        end = new Date(end.getTime() - zoomFactor);
      } else {
        start = new Date(start.getTime() - zoomFactor);
        end = new Date(end.getTime() + zoomFactor);
      }

      Plotly.relayout(chartId, { "xaxis.range": [start, end] });
    },
    { passive: false }
  );

  // Pan
  chartEl.addEventListener("mousedown", function (e) {
    isDragging = true;
    userInteracting = true;
    dragStartX = e.clientX;

    let layout = document.getElementById(chartId).layout;
    dragStartRange = [
      new Date(layout.xaxis.range[0]),
      new Date(layout.xaxis.range[1]),
    ];
  });

  document.addEventListener("mousemove", function (e) {
    if (!isDragging) return;

    let chartWidth = chartEl.offsetWidth;
    let rangeMs = dragStartRange[1] - dragStartRange[0];
    let pxPerMs = chartWidth / rangeMs;
    let deltaMs = (dragStartX - e.clientX) / pxPerMs;

    let newStart = new Date(dragStartRange[0].getTime() + deltaMs);
    let newEnd = new Date(dragStartRange[1].getTime() + deltaMs);

    Plotly.relayout(chartId, { "xaxis.range": [newStart, newEnd] });
  });

  document.addEventListener("mouseup", function () {
    isDragging = false;
  });

  // Start
  loadHistory();
  setInterval(fetchData, 1000);
}
//login funtion
function login() {
  const username = document.getElementById("username").value.trim();
  const password = document.getElementById("password").value.trim();

  fetch("/login", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ username, password }),
  })
    .then((res) => res.json())
    .then((data) => {
      if (data.success) {
        alert("Login successful!");
        // Redirect to dashboard or secured page
        window.location.href = "/dashboard";
      } else {
        alert("Invalid username or password");
      }
    })
    .catch((err) => console.error(err));
}

//logout funtion
function logout() {
  fetch("/logout", { method: "POST" })
    .then((res) => res.json())
    .then((data) => {
      if (data.success) {
        alert("Logged out successfully!");
        window.location.href = "/"; // Redirect to login page
      }
    })
    .catch((err) => console.error(err));
}
//update throshold
function loadThresholds() {
  fetch("/get-thresholds")
    .then((res) => res.json())
    .then((data) => {
      document.getElementById("tempThresh").value = data.temperature;
      document.getElementById("humThresh").value = data.humidity;
      document.getElementById("soundThresh").value = data.sound;
    });
}

function updateThresholds() {
  const thresholds = {
    temperature: parseFloat(document.getElementById("tempThresh").value),
    humidity: parseFloat(document.getElementById("humThresh").value),
    sound: parseFloat(document.getElementById("soundThresh").value),
  };

  fetch("/update-thresholds", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(thresholds),
  })
    .then((res) => res.json())
    .then((data) => {
      if (data.success) alert("Thresholds updated!");
    });
}

//loding the contact and deletion
// Load contacts from CSV
function loadContacts() {
  fetch("/contacts/list")
    .then((res) => res.json())
    .then((data) => {
      const tbody = document.querySelector("#contactsTable tbody");
      tbody.innerHTML = "";
      data.forEach((c) => {
        const row = document.createElement("tr");
        row.innerHTML = `
                <td>${c.name}</td>
                <td>${c.email}</td>
                <td>${c.phone}</td>
                <td><button onclick='deleteContact("${c.name}","${c.email}","${c.phone}")'>Delete</button></td>
            `;
        tbody.appendChild(row);
      });
    });
}

function deleteContact(name, email, phone) {
  if (!confirm(`Delete contact ${name}?`)) return;

  fetch("/contacts/delete", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ name, email, phone }),
  })
    .then((res) => res.json())
    .then((data) => {
      if (data.success) loadContacts();
    });
}

// Add a new contact
function addContact() {
  const contact = {
    name: document.getElementById("name").value.trim(),
    email: document.getElementById("email").value.trim(),
    phone: document.getElementById("phone").value.trim(),
  };

  fetch("/contacts/add", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(contact),
  })
    .then((res) => res.json())
    .then((data) => {
      if (data.success) {
        document.getElementById("name").value = "";
        document.getElementById("email").value = "";
        document.getElementById("phone").value = "";
        loadContacts();
      }
    });
}


    // Initial load
loadContacts();
// Load thresholds when page opens
loadThresholds();
// Initial load
updateTime();
setInterval(updateTime, 1000);
