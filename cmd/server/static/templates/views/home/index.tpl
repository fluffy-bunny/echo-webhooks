{{define "views/home/index"}}
{{template "html_begin" .}}
{{template "header" .}}
{{template "navbar" .}}
<body>
<!-- Page content-->
<div class="container">
    <div class="text-center mt-5">
        <h1>Webhooks Example App with OAuth2 bearer tokens</h1>
        <p class="lead">A complete project boilerplate built with Bootstrap</p>
        <p>Bootstrap v5.1.3</p>
        <div class="alert alert-success" role="alert">
                 <table class="table table-striped">
                <thead>
                <tr>
                <th class="text-start" scope="col">#</th>
                <th class="text-start" scope="col">Type</th>
                <th class="text-start" scope="col">Value</th>
                </tr>
            </thead>
            <tbody>
            {{range $idx,$claim := .claims}}
                <tr>
                <th class="text-start" scope="row">{{$idx}}</th>
                <td class="text-start">{{$claim.Type}}</td>
                <td class="text-start">{{$claim.Value}}</td>
                </tr>
            {{end}}
            </tbody>
            </table>
            <div class="container">
            <div class="columns is-centered is-mobile">	
                <div class="column is-dark notification is-four-fifths">
                    <div class="is-size-7 has-text-warning" id="progress">
                        <h1 class="title">Server Side Events</h1>
                        <h2 class="subtitle">Progress</h2>
                        <ul id="progress_list">
                        </ul>
                    </div>
                </div>
            </div>
        </div>
        <div class="container">
            <div class="columns is-centered is-mobile">	
                <div class="column is-dark notification is-four-fifths">
                    <div class="is-size-7 has-text-warning" id="display">
                        <ul id="display_list">
                        </ul>
                    </div>
                </div>
            </div>
        </div>
        </div>
    </div>
</div>
</body>
<script type="text/javascript">
        var displayList  = document.getElementById("display_list");
        
        e1 = new EventSource('/events/webhooks');
        e1.onmessage = function(event) {
            var li = document.createElement("li");
            li.innerHTML = event.data;
            displayList.appendChild(li);
            var items = document.querySelectorAll("#display_list li");
            if (items.length > 3) {
                displayList.removeChild(items[0]);
            }
        };

        var progressList  = document.getElementById("progress_list");
        

        e2 = new EventSource('/events/webhooks-progress');
        e2.onmessage = function(event) {
            var li = document.createElement("li");
            li.innerHTML = event.data;
            progressList.appendChild(li);
            var items = document.querySelectorAll("#progress_list li");
            if (items.length > 1) {
                progressList.removeChild(items[0]);
            }
 

        };
    </script>

{{template "footer" .}}
{{template "html_end" .}}
{{end}}